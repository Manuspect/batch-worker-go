package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	nc, err := nats.Connect(
		fmt.Sprintf("nats://%s:%s", os.Getenv("NATS_HOST"), os.Getenv("NATS_PORT")),
	)
	if err != nil {
		log.Fatal(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}
	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "BATCH",
		Subjects: []string{"BATCH.>"},
	})
	if err != nil {
		log.Fatal(err)
	}

	cons, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "BATCH",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatal(err)
	}

	cc, err := cons.Consume(
		consumerHandler,
		jetstream.ConsumeErrHandler(
			func(consumeCtx jetstream.ConsumeContext, err error) {
				fmt.Println(err)
			}),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func consumerHandler(msg jetstream.Msg) {
	fmt.Println(string(msg.Data()))
	msg.Ack()
}

// 0. Получить где бетч из очереди
// 1. Получить архив (бэтч) из S3
// 2. Достать из архива видео файл и файл мета-данных
// 3. Преобразовать видео-файл в картинки и сохранить картинки в S3
// 4. Отправить данные из мета-данных (файла) в сервис обработки
