package model

import (
	"context"

	"fmt"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Patient struct {
	MedicalHistory string `json:"medicalhistory" bson:"medicalhistory"`
}

var Client *mongo.Client

func Connection() (*mongo.Client, error) {
	var err error
	const uri = "mongodb+srv://new-user:new-user@cluster0.grve526.mongodb.net/"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = Client.Ping(ctx, nil)
	if err != nil {
		Client.Disconnect(ctx)
		return nil, err
	}

	fmt.Println("Connected to MongoDB successfully!")
	return Client, nil
}

func DisconnectDB(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client.Disconnect(ctx)
}

func FetchData(id string) (Patient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	collection := Client.Database("crud").Collection("patients")
	var patient Patient
	err := collection.FindOne(ctx, bson.M{"patientid": id}).Decode(&patient)
	return patient, err
}
