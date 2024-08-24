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

	log.Infof("register: adding: %v", o)
	db.observers = append(db.observers, o)
}

func (db *Database) unregister(o Observer) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	var idxToRemove int = -1

	for idx, value := range db.observers {
		if value == o {
			idxToRemove = idx
			break
		}
	}

	if idxToRemove >= 0 {
		log.Infof("unregister: removing: %v", o)
		db.observers = append(db.observers[:idxToRemove], db.observers[idxToRemove+1:]...)
	} else {
		log.Errorf("unregister: failed to remove %v", o)
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
	ID   string
	Conn *websocket.Conn
}

type NewRequestMonitor struct {
	Observer
	ID   string
	Conn *websocket.Conn
}

// notify User
func (o *DatabaseMonitor) updated(request *ModRequest) {
	if request.Moderated == true {
		log.Infof("[%s] message %s updated, clientID: %s", o.ID, request.Message.Data, request.ClientID)

		if request.Approved {
			data, err := json.Marshal(request)
			if err != nil {
				// if disconnected, it comes here
				log.Warnf("[%s] json.Marshal failed in updated, %v", o.ID, err)
				return
			}

			err = o.Conn.WriteMessage(websocket.TextMessage, []byte(data))
			if err != nil {
				// if disconnected, it comes here
				log.Warnf("[%s] Mod WriteMessage failed, %v", o.ID, err)
			}
		}
	}
}

// notify Moderator
func (o *NewRequestMonitor) updated(request *ModRequest) {
	if request.UserEmail != "system@system" {
		log.Infof("[%s] new request %s", o.ID, request.Message.Data)
		data, err := json.Marshal(request)
		if err != nil {
			// if disconnected, it comes here
			log.Warnf("[%s] json.Marshal failed in updated, %v", o.ID, err)
			return
		}

		err = o.Conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			// if disconnected, it comes here
			log.Warnf("[%s] Mod WriteMessage failed, %v", o.ID, err)
		}
	}
}
