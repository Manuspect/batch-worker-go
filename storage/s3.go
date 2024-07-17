package storage

import (
	"fmt"
	"os"

	"github.com/gofiber/storage/minio"
	minio_cl "github.com/minio/minio-go/v7"
)

func MinioConnect() *minio_cl.Client {
	store := minio.New(minio.Config{
		Bucket:   os.Getenv("MINIO_BUCKET"),
		Endpoint: fmt.Sprintf("%s:%s", os.Getenv("MINIO_HOST"), os.Getenv("MINIO_SERVER_PORT")),
		Credentials: minio.Credentials{
			AccessKeyID:     os.Getenv("MINIO_USER"),
			SecretAccessKey: os.Getenv("MINIO_PASS"),
		},
	})

	err := store.CheckBucket()
	if err != nil {
		store.CreateBucket()
	}

	return store.Conn()

}
