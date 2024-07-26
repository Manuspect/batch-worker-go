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

type q_msg struct {
	ObjectName string `json:"objectName"`
	BucketName string `json:"bucketName"`
}

func CreateConsumerHandler(m *minio.Client, fc jetstream.Consumer) func(jetstream.Msg) {
	return func(msg jetstream.Msg) {

		var q *q_msg
		json.Unmarshal(msg.Data(), &q)

		ctx := context.Background()
		bucketName := q.BucketName
		objectName := q.ObjectName
		filename := fmt.Sprintf("./%s", objectName)

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

		err = os.RemoveAll("frames")
		if err != nil {
			log.Println(err)
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
	contentType := "application/octet-image"

	arr := strings.Split(objectName, "-")
	var filePath string
	j := 3
	event := events[j]

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
	url := os.Getenv("URL_PROCESSING_SERVICE")
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

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
