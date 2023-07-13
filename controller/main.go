package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"http/model"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Patient struct {
	model.Patient
}

func HandleBulkRequest(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	segments := strings.Split(path, "/")
	id := segments[len(segments)-1]

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	collection := model.Client.Database("crud").Collection("patients")

	if r.Method == http.MethodGet {

		ids := r.URL.Query().Get("ids")
		if ids != "" {
			idList := strings.Split(ids, ",")
			var patients []Patient
			patientChan := make(chan model.Patient)
			done := make(chan struct{})

			for _, patientID := range idList {
				go func(pid string) {
					defer func() {
						done <- struct{}{}
					}()

					patient, err := model.FetchData(pid)
					if err != nil {
						fmt.Printf("Error retrieving details for patient ID %s: %v\n", pid, err)
						patient.MedicalHistory = "Patient not found"
					}
					patientChan <- patient
				}(patientID)
			}

			go func() {
				for range idList {
					<-done
				}
				close(patientChan)
			}()

			for patient := range patientChan {
				patients = append(patients, Patient{Patient: patient})
			}

			responseJSON, err := json.Marshal(patients)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			rw.Header().Set("Content-Type", "application/json")
			rw.Write(responseJSON)
			return
		}
	}

	var patient Patient
	err := collection.FindOne(ctx, bson.M{"patientid": id}).Decode(&patient)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(patient)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(responseJSON)
}
