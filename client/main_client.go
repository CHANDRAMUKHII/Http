package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Patient struct {
	ID string `json:"patientid"`
}

func main() {
	var average time.Duration
	for i := 0; i < 3; i++ {
		startTime := time.Now()
		sendBulkRequest()
		elapsedTime := time.Since(startTime)
		average += elapsedTime
		fmt.Println("Total time taken:", elapsedTime)
	}
	average /= 3
	fmt.Print("Average time taken : ", average)
}

func sendBulkRequest() {

	jsonData, err := ioutil.ReadFile("ids_lakh.json")
	if err != nil {
		log.Fatal(err)
	}

	var patients []Patient

	if err := json.Unmarshal(jsonData, &patients); err != nil {
		log.Fatal(err)
	}

	batchSize := 100
	totalPatients := len(patients)

	for i := 0; i < totalPatients; i += batchSize {
		end := i + batchSize
		if end > totalPatients {
			end = totalPatients
		}

		batch := patients[i:end]
		sendBatchRequest(batch)
	}
}

func sendBatchRequest(batch []Patient) {
	var patientIDs []string
	for _, patient := range batch {
		patientIDs = append(patientIDs, patient.ID)
	}

	getURL := fmt.Sprintf("http://localhost:3000/details?ids=%s", strings.Join(patientIDs, ","))

	response, err := http.Get(getURL)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	fmt.Println(string(responseData))
}
