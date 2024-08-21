package main

import (
	"flag"
	"fmt"
	"net/http"

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
	log.Println("Starting WebSocket Moderator Server...")

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
		// TODO: add the caller as moderator
		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Info(err)
			return
		}
		defer ws.Close()

		moderatorID := uuid.NewString()
		err = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("moderatorID: %s", moderatorID)))
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
		err = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("clientID: %s", clientID)))
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

			// send a message to client.
			msg := "Moderating..."
			log.Infof("[%s] Sending: %s", clientID, msg)
			err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				// if disconnected, it comes here
				log.Warnf("[%s] WriteMessage failed, %v", clientID, err)
				break
			}

			// TODO: moderate clientID and the prompt
			// if moderated, send the prompt to LLM,
			// then moderate, then return the answer or 'moderator refused'
			writeRequest(clientID, string(message), "new")

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
