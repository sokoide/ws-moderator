package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/minio/websocket"
	log "github.com/sirupsen/logrus"
)

type options struct {
	host     string
	port     int
	logLevel string
	level    log.Level
}

var o options = options{
	host:     "localhost",
	port:     80,
	logLevel: "INFO",
}

var clientID string

func parseArgs() {
	flag.StringVar(&o.host, "host", o.host, "Host to talk to")
	flag.IntVar(&o.port, "port", o.port, "Port to talk to")
	flag.StringVar(&o.logLevel, "logLevel", o.logLevel, "Log level")
	flag.Parse()

	level, err := log.ParseLevel(o.logLevel)
	if err == nil {
		o.level = level
	}
}

func connectAndSend() {
	// connect to o.host:o.port and send/receive websocket messages
	u := fmt.Sprintf("ws://%s:%d/chat", o.host, o.port)

	log.Infof("Connecting to: %v", u)

	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Handle interrupt signal to gracefully close the connection
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	// Read messages from the server
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Errorf("Read error: %v", err)
				return
			}
			log.Infof("Received: %s", message)
		}
	}()

	// Send messages to the server
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			msg := "Hello, I'm Scott"
			log.Infof("Sending: %s", msg)
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Errorf("Write error: %v", err)
				return
			}
		case <-interrupt:
			log.Println("Interrupt received, closing connection...")
			// Gracefully close the connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Errorf("Close error: %v", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func main() {
	parseArgs()
	log.SetLevel(o.level)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})

	connectAndSend()
}
