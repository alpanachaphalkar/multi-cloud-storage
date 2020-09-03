package connections

import (
	c "context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"

	"cloud.google.com/go/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

// GetGcpService gets GCP Bucket connection
func GetGcpService() (*storage.BucketHandle, c.Context, error) {
	ctx := c.Background()

	PrivateKeyData, bucketName := os.Getenv("PrivateKeyData"), os.Getenv("bucket_name")
	if len(PrivateKeyData) == 0 || len(bucketName) == 0 {
		log.Fatal("PrivateKeyData or bucket_name environment variable is not set")
	}

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(PrivateKeyData)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "while decoding PrivateKeyData")
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(decodedPrivateKey))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "while establishing gcp storage connection")
	}

	bkt := client.Bucket(bucketName)

	return bkt, ctx, nil
}

// GetAzureService gets Azure Container connection
func GetAzureService() (*azblob.ContainerURL, c.Context, error) {

	accountName, accountKey, containerName := os.Getenv("storageAccountName"), os.Getenv("accessKey"), os.Getenv("containerName")
	if len(accountName) == 0 || len(accountKey) == 0 || len(containerName) == 0 {
		log.Fatal("storageAccountName or accessKey or containerName environment variable is not set")
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "while invalid credentials")
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	containerURL := azblob.NewContainerURL(*URL, p)
	ctx := c.Background()
	return &containerURL, ctx, nil
}
