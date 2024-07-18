package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"processing-worker/internal/batch"
	"strconv"
	"strings"

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

		arr := strings.Split(objectName, "-")
		contentType := "application/octet-image"
		for j := 0; j < len(images); j++ {
			objectName = arr[0] + "-" + arr[1] + "-" + arr[2] + "-" + strconv.Itoa(j+1) + ".jpeg" // 2-3-44-3.jpeg
			filePath := "frames/" + strconv.Itoa(j+1) + ".jpeg"

			_, err := m.FPutObject(
				ctx,
				bucketName,
				objectName,
				filePath,
				minio.PutObjectOptions{ContentType: contentType})
			if err != nil {
				log.Println(err)
			}
		}

		err = os.RemoveAll("frames")
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
