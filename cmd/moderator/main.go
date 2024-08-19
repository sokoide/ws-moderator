package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/minio/websocket"
	log "github.com/sirupsen/logrus"
)

type options struct {
	port     int
	logLevel string
	level    log.Level
}

var o options = options{
	port:     8080,
	logLevel: "INFO",
}

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

	var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

	// "/" for chat
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
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

		log.Infof("[%s] New Connection Established", clientID)

		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				// if disconnected, it comes here
				log.Warnf("[%s] ReadMessage failed, %v", clientID, err)
				break
			}
			log.Infof("[%s] Received Message: %s", clientID, string(message))

			// TODO: moderate clientID and the prompt
			// if moderated, send the prompt to LLM,
			// then moderate, then return the answer or 'moderator refused'

			// send a message to client.
			msg := "hello"
			log.Infof("[%s] Sending: %s", clientID, msg)
			err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				// if disconnected, it comes here
				log.Warnf("[%s] WriteMessage failed, %v", clientID, err)
				break
			}
		}
	})

	// TODO: "/image" for images

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", o.port), nil))
}

func main() {
	parseArgs()
	log.SetLevel(o.level)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})

	startModerator()
}
