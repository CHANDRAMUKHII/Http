package main

import (
	"fmt"
	"http/patient-service/controller"
	"http/patient-service/model"
	"log"
	"net/http"
)

func main() {
	client, _ := model.Connection()
	defer model.DisconnectDB(client)
	http.HandleFunc("/details", controller.HandleBulkRequest)
	log.Fatal(http.ListenAndServe(":8002", nil))
	fmt.Print("Listening on port 8002")
}
