package main

import (
	"fmt"
)

func main() {

	service := handler.getGcpService()
	//fmt.Printf("Starting server at port 8080\n")
	//http.HandleFunc("/storage", azureBlobOperation)
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("GCP service connected")
}
