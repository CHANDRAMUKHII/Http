package model

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"net/http"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Patient struct {
	MedicalHistory string `json:"medicalhistory" bson:"medicalhistory"`
}

var client *mongo.Client

func Connection() (*mongo.Client, error) {
	var err error
	const uri = "mongodb+srv://new-user:new-user@cluster0.grve526.mongodb.net/"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		client.Disconnect(ctx)
		return nil, err
	}

	fmt.Println("Connected to MongoDB successfully!")
	return client, nil
}

func DisconnectDB(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client.Disconnect(ctx)
}

func HandleBulkRequest(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	segments := strings.Split(path, "/")
	id := segments[len(segments)-1]

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	collection := client.Database("crud").Collection("patients")

	if r.Method == http.MethodGet {

		ids := r.URL.Query().Get("ids")
		if ids != "" {
			idList := strings.Split(ids, ",")
			var patients []Patient
			patientChan := make(chan Patient)
			done := make(chan struct{})

			for _, patientID := range idList {
				go func(pid string) {
					defer func() {
						done <- struct{}{}
					}()

					var patient Patient
					err := collection.FindOne(ctx, bson.M{"patientid": pid}).Decode(&patient)
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
				patients = append(patients, patient)
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
