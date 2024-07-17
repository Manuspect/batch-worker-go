package nats

import (
	"context"
	"fmt"
	"os"
	"processing-worker/internal/batch"

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

		objectName := "36.tar.gz"

		bucketName := "dev-bucket"

		ctx := context.Background()
		object, err := m.GetObject(
			ctx,
			bucketName,
			objectName,
			minio.GetObjectOptions{},
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		// object save to file
		// Upload the test file with FPutObject
		// uploadInfo, err := m.FPutObject(ctx, bucketName, filePath, objectName, minio.PutObjectOptions{ContentType: "application/octet-stream"})
		// if err != nil {
		// 	log.Fatalln(err)
		// }

		fmt.Println("OK")
		fmt.Println("================================================")

		images, events, err := batch.GetWebmCsv(objectName)
		if err != nil {
			fmt.Println("GetWebmCsv err:", err)
		}
		if images != nil {
			fmt.Println("images: OK")
		}
		fmt.Println(events[1])
		// 0. Получить где бетч из очереди
		// 1. Получить архив (бэтч) из S3
		// 2. Достать из архива видео файл и файл мета-данных
		// 3. Преобразовать видео-файл в картинки и сохранить картинки в S3
		// 4. Отправить данные из мета-данных (файла) в сервис обработки
		msg.Ack()
	}
}
