package main

import (
	"sync/atomic"
	"time"

	"github.com/minio/websocket"
	log "github.com/sirupsen/logrus"
)

func pingModerator(ws *websocket.Conn, moderatorID string) {
	var err error
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		for {
			select {
			case <-ticker.C:
				log.Infof("tick: %d", atomic.LoadInt32(&claudeConns))
				msg := makeModRequestJsonBytes("", "bot", "", "ping", "ping", true, true)
				err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					// if disconnected, it comes here
					log.Warnf("[%s] Mod WriteMessage failed, %v", moderatorID, err)
					return
				}
			}
		}
	}
}
