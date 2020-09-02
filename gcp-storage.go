package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// BucketObjects Struct for Json response
type BucketObjects struct {
	Bucket      string
	ObjectNames []string
}

func storageOperation(w http.ResponseWriter, r *http.Request) {
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

	bkt := client.Bucket(bucketName)

	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		fmt.Println("Listing the objects in the bucket:")
		query := &storage.Query{Prefix: ""}

		var objectsNames []string
		it := bkt.Objects(ctx, query)
		for {
			attrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			objectsNames = append(objectsNames, attrs.Name)
		}

		bucketObjects := BucketObjects{bucketName, objectsNames}
		encjson, _ := json.Marshal(bucketObjects)
		fmt.Println(string(encjson))
		w.WriteHeader(http.StatusOK)
		w.Write(encjson)
	case "POST":
		mr, err := r.MultipartReader()
		if err != nil {
			log.Fatal(err)
		}
		var fileName string
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				// handle error
				return
			}
			fileName = p.FileName()
			storageW := bkt.Object(fileName).NewWriter(ctx)
			if _, err := io.Copy(storageW, p); err != nil {
				// handle error
				return
			}
			if err := storageW.Close(); err != nil {
				// handle error
				return
			}
		}
		w.Write([]byte(`{"message": "file ` + fileName + ` is uploaded"}`))
	case "DELETE":
		reqFilepath := r.URL.Query().Get("filepath")
		fmt.Printf("Deleting the file %s in the bucket %s\n", reqFilepath, bucketName)
		bkt.Object(reqFilepath).Delete(ctx)
		w.Write([]byte(`{"message": "file ` + reqFilepath + ` is deleted"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}

func main() {
	fmt.Printf("Starting server at port 8080\n")
	http.HandleFunc("/gcpstorage", storageOperation)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
