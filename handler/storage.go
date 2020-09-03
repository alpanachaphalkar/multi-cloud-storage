package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"cloud.google.com/go/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	con "github.com/alpanachaphalkar/multi-cloud-storage/handler/connections"
	"google.golang.org/api/iterator"
)

// Blobs Struct for Json response
type Blobs struct {
	Container string
	Blobs     []string
}

// BucketObjects Struct for Json response
type BucketObjects struct {
	Bucket  string
	Objects []string
}

// GetItems gets the items in the storage Entity
func GetItems(broker string) ([]byte, error) {
	var encjson []byte
	switch broker {
	case "gcp":
		fmt.Println("GCP broker selected")
		bkt, gcpCtx, err := con.GetGcpService()
		if err != nil {
			log.Fatal(err)
		}
		query := &storage.Query{Prefix: ""}
		var objectsNames []string
		it := bkt.Objects(gcpCtx, query)
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
		bucketObjects := BucketObjects{os.Getenv("bucket_name"), objectsNames}
		encjson, _ = json.Marshal(bucketObjects)
		fmt.Println(string(encjson))
	case "azure":
		fmt.Println("Azure broker selected")
		containerURL, azureCtx, err := con.GetAzureService()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Listing the blobs in the container:")
		var blobItems []string
		for marker := (azblob.Marker{}); marker.NotDone(); {
			listBlob, err := containerURL.ListBlobsFlatSegment(azureCtx, marker, azblob.ListBlobsSegmentOptions{})
			if err != nil {
				log.Fatal(err)
			}
			marker = listBlob.NextMarker
			for _, blobInfo := range listBlob.Segment.BlobItems {
				//fmt.Print("	Blob name: " + blobInfo.Name + "\n")
				blobItems = append(blobItems, blobInfo.Name)
			}
		}
		blobs := Blobs{os.Getenv("containerName"), blobItems}
		encjson, _ = json.Marshal(blobs)
		fmt.Println(string(encjson))
	default:
		encjson = []byte(`{"message": "No valid broker selected"}`)
	}
	return encjson, nil
}

// UploadFile uploads the file to Storage Entity
func UploadFile(broker string, p *multipart.Part) ([]byte, error) {
	var encjson []byte
	switch broker {
	case "gcp":
		fmt.Println("GCP broker selected")
		bkt, gcpCtx, err := con.GetGcpService()
		if err != nil {
			log.Fatal(err)
		}
		storageW := bkt.Object(p.FileName()).NewWriter(gcpCtx)
		if _, err := io.Copy(storageW, p); err != nil {
			log.Fatal(err)
		}
		if err := storageW.Close(); err != nil {
			log.Fatal(err)
		}
		encjson = []byte(`{"message": "file ` + p.FileName() + ` is uploaded"}`)
		fmt.Println(string(encjson))
	case "azure":
		fmt.Println("Azure broker selected")
		containerURL, azureCtx, err := con.GetAzureService()
		if err != nil {
			log.Fatal(err)
		}
		blobURL := containerURL.NewBlockBlobURL(p.FileName())
		bufferSize := 2 * 1024 * 1024
		maxBuffers := 3
		_, err = azblob.UploadStreamToBlockBlob(azureCtx, p, blobURL, azblob.UploadStreamToBlockBlobOptions{
			BufferSize: bufferSize, MaxBuffers: maxBuffers})
		if err != nil {
			log.Fatal(err)
		}
		encjson = []byte(`{"message": "file ` + p.FileName() + ` is uploaded"}`)
		fmt.Println(string(encjson))
	default:
		encjson = []byte(`{"message": "No valid broker selected"}`)
	}
	return encjson, nil
}

// DeleteItem deletes the file from Storage Entity
func DeleteItem(broker string, filepath string) ([]byte, error) {
	var encjson []byte
	switch broker {
	case "gcp":
		fmt.Println("GCP broker selected")
		bkt, gcpCtx, err := con.GetGcpService()
		if err != nil {
			log.Fatal(err)
		}
		bkt.Object(filepath).Delete(gcpCtx)
		encjson = []byte(`{"message": "file ` + filepath + ` is deleted"}`)
		fmt.Println(string(encjson))
	case "azure":
		fmt.Println("Azure broker selected")
		containerURL, azureCtx, err := con.GetAzureService()
		blobURL := containerURL.NewBlockBlobURL(filepath)
		blobURL.Delete(azureCtx, "include", azblob.BlobAccessConditions{})
		if err != nil {
			log.Fatal(err)
		}
		encjson = []byte(`{"message": "file ` + filepath + ` is deleted"}`)
		fmt.Println(string(encjson))
	default:
		encjson = []byte(`{"message": "No valid broker selected"}`)
	}
	return encjson, nil
}
