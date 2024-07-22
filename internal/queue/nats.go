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

		images, events, err := batch.GetWebmCsv(filename)
		if err != nil {
			fmt.Println("GetWebmCsv err:", err)
		}
		if images != nil {
			fmt.Println("images: OK")
		}
		if events != nil {
			fmt.Println("events: OK")
		}

		type Info struct {
			Timestamp     string `json:"timestamp"`
			Process_path  string `json:"process_path"`
			Title         string `json:"title"`
			Class_name    string `json:"class_name"`
			Window_left   string `json:"window_left"`
			Window_top    string `json:"window_top"`
			Window_right  string `json:"window_right"`
			Window_bottom string `json:"window_bottom"`
			Event         string `json:"event"`
			Mouse_x_pos   string `json:"mouse_x_pos"`
			Mouse_y_pos   string `json:"mouse_y_pos"`
			Modifiers     string `json:"modifiers"`
			FrameName     string `json:"frameName"`
		}

		arr := strings.Split(objectName, "-")
		contentType := "application/octet-image"
		var frame_name string

		for j, event := range events {

			if j > 0 {

				arr_csv := strings.Split(fmt.Sprintf("%v", event), "|")

				frame_name = arr[0] + "-" + arr[1] + "-" + arr_csv[0][1:] + "-" + strconv.Itoa(j) + ".jpeg" // 2-3-44-3.jpeg
				filePath := "frames/" + strconv.Itoa(j) + ".jpeg"

				_, err := m.FPutObject(
					ctx,
					bucketName,
					frame_name,
					filePath,
					minio.PutObjectOptions{ContentType: contentType})
				if err != nil {
					log.Println(err)
				}

				fileInfo := Info{
					Timestamp:     arr_csv[0][1:],
					Process_path:  arr_csv[1],
					Title:         arr_csv[2],
					Class_name:    arr_csv[3],
					Window_left:   arr_csv[4],
					Window_top:    arr_csv[5],
					Window_right:  arr_csv[6],
					Window_bottom: arr_csv[7],
					Event:         arr_csv[8],
					Mouse_x_pos:   arr_csv[9],
					Mouse_y_pos:   arr_csv[10],
					Modifiers:     arr_csv[11][0 : len(arr_csv[11])-1],
					FrameName:     frame_name,
				}

				body, err := json.Marshal(fileInfo)
				if err != nil {
					log.Println(err)
				}

				fmt.Println("fileInfo:", fileInfo)

				bodyReader := bytes.NewReader(body)

				url := "http://localhost:3000/info" // TODO
				req, err := http.NewRequest("Get", url, bodyReader)
				if err != nil {
					log.Println(err)
				}

				req.Header.Set("Content-Type", "application/json")

				client := http.Client{
					Timeout: 5 * time.Second,
				}

				resp, err := client.Do(req)
				if err != nil {
					log.Println(err)
				}

				b, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
				}
				fmt.Println(string(b))

			}
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

// 0. Получить где бетч из очереди
// 1. Получить архив (бэтч) из S3
// 2. Достать из архива видео файл и файл мета-данных
// 3. Преобразовать видео-файл в картинки и сохранить картинки в S3
// 4. Отправить данные из мета-данных (файла) в сервис обработки
