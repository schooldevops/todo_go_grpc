package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "mydb"
const mongoURI = "mongodb://localhost"
const mongoPort = "27017"

func NewMongoCollection(collectionName string) *mongo.Collection {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Printf("New Mongo Collection: %v", collectionName)
	
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI + ":" + mongoPort))
	if err != nil {
		log.Fatalf("Cannot create client for mongodb %v", err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatalf("Couldn't connect to mongodb %v", err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	return collection
}

