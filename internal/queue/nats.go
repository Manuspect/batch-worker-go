package nats

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"processing-worker/internal/batch"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func NatsConnect() (jetstream.JetStream, error) {
	nc, err := nats.Connect(
		fmt.Sprintf("nats://%s:%s", os.Getenv("NATS_HOST"), os.Getenv("NATS_PORT")),
	)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func CreateBatchComsumer(js jetstream.JetStream) (jetstream.Consumer, error) {
	ctx := context.Background()
	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "BATCH",
		Subjects: []string{"BATCH.>"},
	})
	if err != nil {
		return nil, err
	}

	consumer, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "BATCH",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
func CreateFileComsumer(js jetstream.JetStream) (jetstream.Consumer, error) {
	ctx := context.Background()
	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "FILE",
		Subjects: []string{"FILE.>"},
	})
	if err != nil {
		return nil, err
	}

	consumer, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "FILE",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func CreateConsumerHandler(m *minio.Client, fc jetstream.Consumer) func(jetstream.Msg) {
	return func(msg jetstream.Msg) {
		fmt.Println(string(msg.Data()))
		// msg.Ack()

		type q_msg struct {
			ObjectName string `json:"objectName"`
			BucketName string `json:"bucketName"`
		}

		var q *q_msg
		json.Unmarshal(msg.Data(), &q)

		objectName := q.ObjectName

		bucketName := q.BucketName

		filename := fmt.Sprintf("./%s", objectName)

		ctx := context.Background()
		err := m.FGetObject(
			ctx,
			bucketName,
			objectName,
			filename,
			minio.GetObjectOptions{},
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		events, err := batch.GetWebmCsv(filename)
		if err != nil {
			fmt.Println("GetWebmCsv err:", err)
			return
		}

		err = sendJpegCsvFiles(m, objectName, bucketName, events)
		if err != nil {
			fmt.Println("sendJpegCsvFiles err:", err)
			return
		}

		err = os.RemoveAll("frames")
		if err != nil {
			log.Println(err)
		}

		err = os.RemoveAll(objectName)
		if err != nil {
			log.Println(err)
		}

		msg.Ack()
	}
}

func sendJpegCsvFiles(m *minio.Client, objectName, bucketName string, events []batch.FileCsv) error {
	ctx := context.Background()
	arr := strings.Split(objectName, "-")
	contentType := "application/octet-image"

	var filePath string
	for j, event := range events {

		if j > 0 {

			filePath = arr[0] + "-" + arr[1] + "-" + event.Timestamp + "-" + strconv.Itoa(j) + ".jpeg" // 2-3-44-3.jpeg
			filePathJpeg := "frames/" + strconv.Itoa(j) + ".jpeg"

			_, err := m.FPutObject(
				ctx,
				bucketName,
				filePath,
				filePathJpeg,
				minio.PutObjectOptions{ContentType: contentType})
			if err != nil {
				log.Println(err)
				return err
			}

			fileInfo := batch.Info{
				Timestamp:     event.Timestamp,
				Process_path:  event.Process_path,
				Title:         event.Title,
				Class_name:    event.Class_name,
				Window_left:   event.Window_left,
				Window_top:    event.Window_top,
				Window_right:  event.Window_right,
				Window_bottom: event.Window_bottom,
				Event:         event.Event,
				Mouse_x_pos:   event.Mouse_x_pos,
				Mouse_y_pos:   event.Mouse_y_pos,
				Modifiers:     event.Modifiers,
				FilePath:      filePath,
			}

			body, err := json.Marshal(fileInfo)
			if err != nil {
				return err
			}

			bodyReader := bytes.NewReader(body)
			url := "http://localhost:3000/api_v1/info"
			req, err := http.NewRequest(http.MethodPost, url, bodyReader)
			if err != nil {
				return err
			}

			req.Header.Add("Content-Type", "application/json")
			client := http.Client{
				Timeout: 5 * time.Second,
			}

			resp, err := client.Do(req)
			if err != nil {
				return err
			}

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			fmt.Println(string(b))
		}
	}

	return nil
}

// 0. Получить где бетч из очереди
// 1. Получить архив (бэтч) из S3
// 2. Достать из архива видео файл и файл мета-данных
// 3. Преобразовать видео-файл в картинки и сохранить картинки в S3
// 4. Отправить данные из мета-данных (файла) в сервис обработки
