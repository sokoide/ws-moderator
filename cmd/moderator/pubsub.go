package main

import (
	"encoding/json"
	"sync"

	"github.com/minio/websocket"
	log "github.com/sirupsen/logrus"
)

// Subject
type Subject interface {
	register(o Observer)
	unregiter(o Observer)
	notifyObservers(request *ModRequest)
}

type Database struct {
	observers []Observer
	mtx       sync.Mutex
}

func (db *Database) updated(request *ModRequest) {
	db.notifyObservers(request)
}

func (db *Database) register(o Observer) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	db.observers = append(db.observers, o)
}

func (db *Database) unregister(o Observer) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	var idxToRemove int = -1

	for idx, value := range db.observers {
		if value == o {
			idxToRemove = idx
		}
	}

	if idxToRemove >= 0 {
		db.observers = append(db.observers[:idxToRemove], db.observers[idxToRemove+1:]...)
	}
}

func (db *Database) notifyObservers(request *ModRequest) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	for _, o := range db.observers {
		o.updated(request)
	}
}

// Observer
type Observer interface {
	updated(request *ModRequest)
}

type DatabaseMonitor struct {
	Observer
	ID string
}

type NewRequestMonitor struct {
	Observer
	ID   string
	Conn *websocket.Conn
}

func (o *DatabaseMonitor) updated(request *ModRequest) {
	if request.Kind == "updated" {
		log.Infof("[%s] message %s updated, clientID: %s", o.ID, request.Message.Data, request.ClientID)

		// 		log.Infof("[%s] Mod Sending: %s", moderatorID, msg)
		// 		err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
		// 		if err != nil {
		// 			// if disconnected, it comes here
		// 			log.Warnf("[%s] Mod WriteMessage failed, %v", moderatorID, err)
		// 			break
		// 		}
	}
}

func (o *NewRequestMonitor) updated(request *ModRequest) {
	if request.Kind == "new" {
		log.Infof("[%s] new request %s", o.ID, request.Message.Data)
		data, err := json.Marshal(request)
		if err != nil {
			// if disconnected, it comes here
			log.Warnf("[%s] Mod WriteMessage failed in json.Marshal, %v", o.ID, err)
			return
		}

		err = o.Conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			// if disconnected, it comes here
			log.Warnf("[%s] Mod WriteMessage failed, %v", o.ID, err)
		}
	}
}
