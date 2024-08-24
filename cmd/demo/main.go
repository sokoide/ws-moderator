package main

import (
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sokoide/ws-ai/pkg/claude"
)

var claudeConns int32

var claudes sync.Map

type ClaudeComm struct {
	user string
	cin  chan claude.Request
	cout chan claude.Response
}

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

func query(c *ClaudeComm, q string) {
	atomic.AddInt32(&claudeConns, 1)
	user := c.user

	// 1st q
	c.cin <- claude.Request{Prompt: q}
	res := <-c.cout
	log.Infof("[%s] %+v", user, res)
	// 2nd q
	c.cin <- claude.Request{Prompt: "Summarize it in 1 line"}
	res = <-c.cout
	log.Infof("[%s] %+v", user, res)

	atomic.AddInt32(&claudeConns, -1)
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
	// go startUser("foo@localhost", "make a story of a dragon in 100 words")
	// go startUser("bar@localhost", "what's the biggest city in India?")
	// go startUser("baz@localhost", "what's the largest animal on earth?")

	for i := 0; i < 15; i++ {
		log.Infof("%d conns\n", atomic.LoadInt32(&claudeConns))
		time.Sleep(time.Second)
	}

	log.Info("done")
}
