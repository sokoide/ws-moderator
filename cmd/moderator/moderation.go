package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
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

func newModRequest(id string, clientID string, userEmail string, kind string, text string) *ModRequest {
	return &ModRequest{
		ID:        id,
		ClientID:  clientID,
		User:      "system",
		UserEmail: userEmail,
		Message:   Message{Kind: kind, Data: text},
		Approved:  false,
		Moderated: false,
	}
}

func makeModRequestJsonBytes(id string, clientID string, userEmail string, kind string, text string) []byte {
	m := newModRequest(id, clientID, userEmail, kind, text)
	mj, err := json.Marshal(m)

	if err != nil {
		log.Info(err)
		return make([]byte, 0)
	}

	return mj
}
