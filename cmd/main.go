package main

import (
	"log"
	"os"
	"os/signal"
	nats "processing-worker/internal/queue"
	"processing-worker/storage"
	"syscall"

	"fmt"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	godotenv.Load()
	jetStream, err := nats.NatsConnect()
	if err != nil {
		log.Fatalln(err)
	}

	minio_client := storage.MinioConnect()

	bc, err := nats.CreateBatchComsumer(jetStream)
	if err != nil {
		log.Fatal(err)
	}

	fc, err := nats.CreateFileComsumer(jetStream)
	if err != nil {
		log.Fatal(err)
	}

	cc, err := bc.Consume(
		// Debili hodiyt suda
		nats.CreateConsumerHandler(minio_client, fc),
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
