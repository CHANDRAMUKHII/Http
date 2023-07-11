package main

import (
	"fmt"
	"http/model"
	"log"
	"net/http"
)

func main() {
	client, _ := model.Connection()
	defer model.DisconnectDB(client)

	http.HandleFunc("/details", model.HandleBulkRequest)
	log.Fatal(http.ListenAndServe(":3000", nil))
	fmt.Print("Listening in port 3000")
}
