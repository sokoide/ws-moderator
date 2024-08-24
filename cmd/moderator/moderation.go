package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Kind string `json:"kind"` // txt, png or jpg
	Data string `json:"data"` // txt or base64 encoded image
}

type ModRequest struct {
	ID        string  `json:"id"`        // MongoDB unique ID
	ClientID  string  `json:"client_id"` // UUID
	User      string  `json:"user"`
	UserEmail string  `json:"user_email"`
	Message   Message `json:"message"`
	Approved  bool    `json:"approved"`
	Moderated bool    `json:"moderated"`
}

type MongoRequest struct {
	ID          primitive.ObjectID `bson:"_id"`
	ClientID    string             `bson:"client_id"`
	UserEmail   string             `bson:"user_email"`
	MessageData string             `bson:"message_data"`
	MessageKind string             `bson:"message_kind"`
	Approved    bool               `bson:"approved"`
	Moderated   bool               `bson:"moderated"`
}

func newModRequest(id string, clientID string, userEmail string, kind string, text string, approved bool, moderated bool) *ModRequest {
	return &ModRequest{
		ID:        id,
		ClientID:  clientID,
		UserEmail: userEmail,
		Message:   Message{Kind: kind, Data: text},
		Approved:  approved,
		Moderated: moderated,
	}
}

func makeModRequestJsonBytes(id string, clientID string, userEmail string, kind string, text string, approved bool, moderated bool) []byte {
	m := newModRequest(id, clientID, userEmail, kind, text, approved, moderated)
	mj, err := json.Marshal(m)

	if err != nil {
		log.Info(err)
		return make([]byte, 0)
	}

	return mj
}
