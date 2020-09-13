package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/xPlorinRolyPoly/multi-cloud-storage/handler"
)

func main() {

	fmt.Printf("Starting server at port 8080\n")
	http.HandleFunc("/storage", storageOperation)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

func storageOperation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqBroker := r.URL.Query().Get("broker")
	switch r.Method {
	case "GET":
		encjson, err := handler.GetItems(reqBroker)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(encjson)
	case "POST":
		mr, err := r.MultipartReader()
		if err != nil {
			log.Fatal(err)
		}
		var encjson []byte
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			encjson, err = handler.UploadFile(reqBroker, p)
			if err != nil {
				log.Fatal(err)
			}
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(encjson)
	case "DELETE":
		reqFilepath := r.URL.Query().Get("filepath")
		encjson, err := handler.DeleteItem(reqBroker, reqFilepath)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(encjson)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}
