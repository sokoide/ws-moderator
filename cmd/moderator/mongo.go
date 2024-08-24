package main

import (
	"context"

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

func storeRequest(clientID string, userEmail string, messageData string, messageKind string, approved bool, moderated bool) string {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	document := map[string]interface{}{
		"client_id":    clientID,
		"user_email":   userEmail,
		"message_data": messageData,
		"message_kind": messageKind,
		"approved":     approved,
		"moderated":    moderated,
	}

	request := newModRequest("", clientID, userEmail, messageKind, messageData, approved, moderated)
	collection := client.Database(MONGODB_NAME).Collection(MONGODB_COLLECTION)

	insertResult, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Errorf("Failed to write %v in MongoDB, %v", document, err)
		return ""
	}

	log.Infof("Inserted a single document: %v\n", insertResult.InsertedID)

	id, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Errorf("Failed to get MongoDB ObjectID for %s", document)
		return ""
	}

	request.ID = id.Hex()
	db.notifyObservers(request)
	return request.ID
}

func updateRequest(msgid string, approved bool, moderated bool) *ModRequest {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	collection := client.Database(MONGODB_NAME).Collection(MONGODB_COLLECTION)
	// update document whose _id is msgid
	objectID, err := primitive.ObjectIDFromHex(msgid)
	filter := bson.D{{"_id", objectID}}

	request := bson.D{
		{"$set", bson.D{
			{"approved", approved},
			{"moderated", moderated},
		}},
	}
	updated, err := collection.UpdateOne(ctx, filter, request)
	log.Debugf("updated: %v", updated)

	if err != nil {
		log.Errorf("Failed to update %s, %v", msgid, err)
		return nil
	}

	// get the doc
	var result MongoRequest
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Errorf("Failed to get the updated document %s, %v", msgid, err)
		return nil
	}
	log.Debugf("Document found: %v", result)

	modRequest := newModRequest(result.ID.Hex(), result.ClientID, result.UserEmail, result.MessageKind, result.MessageData, result.Approved, result.Moderated)
	db.notifyObservers(modRequest)
	return modRequest
}

func loadRequestsForUserEmail(userEmail string) []ModRequest {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	requests := make([]ModRequest, 0)

	collection := client.Database(MONGODB_NAME).Collection(MONGODB_COLLECTION)
	// get all documents whose user_email == userEmail and approved
	filter := bson.M{"user_email": userEmail, "approved": true}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Errorf("Failed to get MongoDB Objects for user_email == %s", userEmail)
		return requests
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var document map[string]interface{}
		if err := cursor.Decode(&document); err != nil {
			log.Errorf("error %v in loadRequestsForUserEmail cursor.Decode", err)
		} else {
			log.Debugf("Found document: %v\n", document)

			request := newModRequest(document["_id"].(primitive.ObjectID).Hex(),
				document["client_id"].(string),
				document["user_email"].(string),
				document["message_kind"].(string),
				document["message_data"].(string),
				document["approved"].(bool),
				document["moderated"].(bool))
			requests = append(requests, *request)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Errorf("error %v in loadRequestsForUserEmail", err)
	}
	return requests
}

func loadRequests(moderated bool) []ModRequest {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	requests := make([]ModRequest, 0)

	collection := client.Database(MONGODB_NAME).Collection(MONGODB_COLLECTION)
	// get all documents whose approved==false && moderated==moderated
	filter := bson.M{"moderated": moderated, "approved": false}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Errorf("Failed to get MongoDB Objects for moderated=%v", moderated)
		return requests
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var document map[string]interface{}
		if err := cursor.Decode(&document); err != nil {
			log.Errorf("error %v in loadRequests cursor.Decode", err)
		} else {
			log.Debugf("Found document: %v\n", document)

			request := newModRequest(document["_id"].(primitive.ObjectID).Hex(),
				document["client_id"].(string),
				document["user_email"].(string),
				document["message_kind"].(string),
				document["message_data"].(string),
				document["approved"].(bool),
				document["moderated"].(bool))
			requests = append(requests, *request)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Errorf("error %v in loadRequestsl", err)
	}
	return requests
}
