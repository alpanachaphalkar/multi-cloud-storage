package handler

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func getService() *storage.Client {
	ctx := context.Background()

	PrivateKeyData, bucketName := os.Getenv("PrivateKeyData"), os.Getenv("bucket_name")
	if len(PrivateKeyData) == 0 || len(bucketName) == 0 {
		log.Fatal("PrivateKeyData or bucket_name environment variable is not set")
	}

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(PrivateKeyData)
	if err != nil {
		log.Fatal(err)
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(decodedPrivateKey))
	if err != nil {
		log.Fatal(err)
	}
	return client
}
