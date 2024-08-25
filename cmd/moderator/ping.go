package main

import (
	"time"

	"github.com/minio/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/sokoide/ws-ai/pkg/claude"
	"github.com/sokoide/ws-ai/pkg/dalle"
)

func pingModerator(ws *websocket.Conn, moderatorID string) {
	var err error
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		for {
			select {
			case <-ticker.C:
				log.Infof("tick: current claude conns: %d, dalle conns: %d", claude.GetConns(), dalle.GetConns())
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
