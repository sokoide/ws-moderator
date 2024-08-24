package main

import (
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sokoide/ws-ai/pkg/claude"
)

var claudeConns int32

func startUser(id string, firstQ string) {
	cin := make(chan claude.Request, 1)
	cout := make(chan claude.Response, 1)

	go claude.StartConversation(id, cin, cout)
	atomic.AddInt32(&claudeConns, 1)

	for _, q := range []string{firstQ, "Summarize it in 1 sentence"} {
		cin <- claude.Request{Prompt: q}
		res := <-cout

		if res.Succeeded {
			log.Infof("[%s] ans: %s\n", id, res.Text)
		} else {
			log.Infof("[%s] error: %s\n", id, res.Text)
		}
	}
	close(cin)
	atomic.AddInt32(&claudeConns, -1)
}

func main() {
	log.Info("demo")

	go startUser("foo@localhost", "make a story of a dragon in 100 words")
	go startUser("bar@localhost", "what's the biggest city in India?")
	go startUser("baz@localhost", "what's the largest animal on earth?")

	for i := 0; i < 15; i++ {
		log.Infof("%d conns\n", atomic.LoadInt32(&claudeConns))
		time.Sleep(time.Second)
	}

	log.Info("done")
}
