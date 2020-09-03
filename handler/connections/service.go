package connections

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

func getGcpService() (*storage.Client, error) {
	ctx := context.Background()

	PrivateKeyData, bucketName := os.Getenv("PrivateKeyData"), os.Getenv("bucket_name")
	if len(PrivateKeyData) == 0 || len(bucketName) == 0 {
		log.Fatal("PrivateKeyData or bucket_name environment variable is not set")
	}

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(PrivateKeyData)
	if err != nil {
		return nil, errors.Wrapf(err, "while decoding PrivateKeyData")
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(decodedPrivateKey))
	if err != nil {
		return nil, errors.Wrapf(err, "while establishing gcp storage connection")
	}
	return client, nil
}
