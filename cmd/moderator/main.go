package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/websocket"
	log "github.com/sirupsen/logrus"
)

// globals
type options struct {
	port     int
	logLevel string
	level    log.Level
}

var o options = options{
	port:     80,
	logLevel: "INFO",
}

var db Database = Database{}

// functions
func parseArgs() {
	flag.IntVar(&o.port, "port", o.port, "Port to listen on")
	flag.StringVar(&o.logLevel, "logLevel", o.logLevel, "Log level")
	flag.Parse()

	level, err := log.ParseLevel(o.logLevel)
	if err == nil {
		o.level = level
	}
}

func startModerator() {
	log.Info("Starting WebSocket Moderator Server...")

	var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections by returning true
			return true
			// Alternatively,
			// return r.Header.Get("Origin") == "http://allowed-origin.com"
		}}

	// "/" for react

	// "/moderator" from clients
	http.HandleFunc("/moderator", func(writer http.ResponseWriter, request *http.Request) {
		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Info(err)
			return
		}
		defer ws.Close()

		moderatorID := uuid.NewString()
		msg := makeModRequestJsonBytes("", "bot", "system@system", "txt", fmt.Sprintf("moderatorID: %s", moderatorID), true, true)
		err = ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Error(err)
			return
		}

		monitor := &NewRequestMonitor{
			ID:   moderatorID,
			Conn: ws,
		}
		db.register(monitor)
		defer db.unregister(monitor)

		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()

		for {
			// _, message, err := ws.ReadMessage()
			// if err != nil {
			// 	// if disconnected, it comes here
			// 	log.Warnf("[%s] ReadMessage failed, %v", moderatorID, err)
			// 	break
			// }
			// log.Infof("[%s] Received Message: %s", moderatorID, string(message))

			// TODO: moderate clientID and the prompt
			// if moderated, send the prompt to LLM,
			// then moderate, then return the answer or 'moderator refused'

			// send a message to the moderator service
			// msg := "hello"
			// log.Infof("[%s] Sending: %s", moderatorID, msg)
			// err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
			// if err != nil {
			// 	// if disconnected, it comes here
			// 	log.Warnf("[%s] WriteMessage failed, %v", moderatorID, err)
			// 	break
			// }

			// for {
			// 	select {
			// 	case req := <-modQ:
			// 		log.Infof("[%s] Mod req: %+v", moderatorID, req)
			// 		msg := fmt.Sprintf("[%s] req: %s", req.ClientID, req.Message.Data)

			// 		log.Infof("[%s] Mod Sending: %s", moderatorID, msg)
			// 		err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
			// 		if err != nil {
			// 			// if disconnected, it comes here
			// 			log.Warnf("[%s] Mod WriteMessage failed, %v", moderatorID, err)
			// 			break
			// 		}

			// 		req.ReturnChannel <- true
			// 	}
			// }
			for {
				select {
				case <-ticker.C:
					log.Debug("tick")
					msg := makeModRequestJsonBytes("", moderatorID, "system@system", "ping", "ping", true, true)
					err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
					if err != nil {
						// if disconnected, it comes here
						log.Warnf("[%s] Mod WriteMessage failed, %v", moderatorID, err)
						return
					}
				}
			}
		}

	})

	// "/chat" from clients
	http.HandleFunc("/chat", func(writer http.ResponseWriter, request *http.Request) {
		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Info(err)
			return
		}
		defer ws.Close()

		clientID := uuid.NewString()
		msg := makeModRequestJsonBytes("", "bot", "system@system", "txt", fmt.Sprintf("clientID: %s", clientID), true, true)
		err = ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Error(err)
			return
		}

		monitor := &DatabaseMonitor{
			ID: clientID,
		}
		db.register(monitor)
		defer db.unregister(monitor)

		log.Infof("[%s] New Connection Established", clientID)

		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				// if disconnected, it comes here
				log.Warnf("[%s] ReadMessage failed, %v", clientID, err)
				break
			}
			log.Infof("[%s] Received Message: %s", clientID, string(message))

			// convert message into ModRequest
			var req ModRequest
			err = json.Unmarshal(message, &req)
			if err != nil {
				log.Errorf("[%s] failed to parse %s", clientID, string(message))
				continue
			}
			// save it
			storeRequest(clientID, req.UserEmail, req.Message.Data, req.Message.Kind, false, false)

			// send a message to client.
			msg := makeModRequestJsonBytes("", clientID, "system@system", "txt", "Moderating...", true, true)
			err = ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				// if disconnected, it comes here
				log.Warnf("[%s] WriteMessage failed, %v", clientID, err)
				break
			}

			// call AI
			if strings.HasPrefix(req.Message.Data, "/imagine") {
				// TODO: image generation
			} else {
				// TODO: text generation
				// time.Sleep(time.Second)

				// storeRequest(clientID, req.UserEmail, "Dummy response from Claude3...", "txt", false, false)
				storeRequest("bot", req.UserEmail, "Dummy response from Claude3...", "txt", true, true)
				// msg := makeModRequestJsonBytes("", clientID, "system@system", "txt", "Dummy answer from Claude3...")
				// err = ws.WriteMessage(websocket.TextMessage, msg)
				// if err != nil {
				// 	// if disconnected, it comes here
				// 	log.Warnf("[%s] WriteMessage failed, %v", clientID, err)
				// 	break
				// }
			}

			// TODO: moderate clientID and the prompt
			// if moderated, send the prompt to LLM,
			// then moderate, then return the answer or 'moderator refused'

		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", o.port), nil))
}

// main
func main() {
	parseArgs()
	log.SetLevel(o.level)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})

	startModerator()
}
