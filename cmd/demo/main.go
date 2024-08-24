package main

import (
	"math/rand"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sokoide/ws-ai/pkg/claude"
)

var claudes sync.Map

type ClaudeComm struct {
	user string
	cin  chan claude.Request
	cout chan claude.Response
}

func query(c *ClaudeComm, q string) {
	user := c.user

	// 1st q
	c.cin <- claude.Request{Prompt: q}
	res := <-c.cout
	log.Infof("[%s] %+v", user, res)

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))
	// 2nd q
	c.cin <- claude.Request{Prompt: "Summarize it in 1 line"}
	res = <-c.cout
	log.Infof("[%s] %+v", user, res)

	log.Infof("Closing %s\n", c.user)
	close(c.cin)
}

func main() {
	log.Info("demo")

	users := []string{"foo@localhost", "bar@localhost", "baz@localhost"}
	questions := []string{"What's the biggest city in Japan?", "What's the biggest animal?", "Make a story of a dragon in 100 words."}

	for _, user := range users {
		if value, ok := claudes.Load(user); ok {
			log.Infof("%s->%v", user, value)
		} else {
			log.Infof("%s not available, making...", user)
			c := &ClaudeComm{
				user: user,
				cin:  make(chan claude.Request, 1),
				cout: make(chan claude.Response, 1),
			}
			claudes.Store(user, c)
			go claude.StartConversation(c.user, c.cin, c.cout)
		}
	}

	for idx, user := range users {
		q := questions[idx]
		if c, ok := claudes.Load(user); ok {
			go query(c.(*ClaudeComm), q)
		}
	}

	for i := 0; i < 15; i++ {
		log.Infof("%d conns\n", claude.GetConns())
		time.Sleep(time.Second)
	}

	log.Info("done")
}
