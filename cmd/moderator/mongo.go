package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

// global
var client *mongo.Client

const MONGODB_NAME = "familyday"
const MONGODB_COLLECTION = "requests"

func init() {
	var err error

	ctx := context.TODO()
	clientOptions := mongoOptions.Client().ApplyURI("mongodb://root:password@localhost:27017/?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
}

func writeRequest(clientID string, messageText string, messageType string) {
	ctx := context.TODO()

	document := map[string]interface{}{
		"clientID":    clientID,
		"messageText": messageText,
		"messageType": messageType,
		"moderated":   false,
		"approved":    false,
	}

	collection := client.Database(MONGODB_NAME).Collection(MONGODB_COLLECTION)

	insertResult, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Errorf("Failed to write in MongoDB, %v", err)
		return
	}

	log.Infof("Inserted a single document: %v\n", insertResult.InsertedID)

	id, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Errorf("Failed to get MongoDB ObjectID for %s", document)
		return
	}

	request := &ModRequest{
		ID:       id.Hex(),
		ClientID: clientID,
		Kind:     messageType,
		Message: Message{
			Kind: "txt",
			Data: messageText,
		},
	}

	db.notifyObservers(request)
}

func writeMongoSpike() {
	// Set up a context and options
	ctx := context.TODO()
	clientOptions := mongoOptions.Client().ApplyURI("mongodb://root:password@localhost:27017/?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure disconnection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Select the database and collection
	collection := client.Database("familyday").Collection("requests")

	// Insert a single document
	document := map[string]interface{}{
		"name":  "John Doe",
		"email": "johndoe@example.com",
		"age":   30,
	}

	insertResult, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Inserted a single document: %v\n", insertResult.InsertedID)

	// Insert multiple documents
	documents := []interface{}{
		map[string]interface{}{
			"name":  "Alice",
			"email": "alice@example.com",
			"age":   28,
		},
		map[string]interface{}{
			"name":  "Bob",
			"email": "bob@example.com",
			"age":   35,
		},
	}

	insertManyResult, err := collection.InsertMany(ctx, documents)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Inserted multiple documents: %v\n", insertManyResult.InsertedIDs)

	// query
	var result bson.M
	err = collection.FindOne(ctx, bson.M{"name": "John Doe"}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)

	// Query for multiple documents
	cursor, err := collection.Find(ctx, bson.M{"age": bson.M{"$gt": 30}})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found multiple documents: %+v\n", results)

	// Query with options (e.g., sorting and limiting)
	findOptions := mongoOptions.Find()
	findOptions.SetSort(bson.D{{"age", 1}})
	findOptions.SetLimit(2)

	cursor, err = collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	results = nil
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found documents with options: %+v\n", results)
}
