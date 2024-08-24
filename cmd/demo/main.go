package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sokoide/ws-ai/pkg/claude"
)

func main() {
	log.Info("demo")

	cin := make(chan claude.Request, 1)
	cout := make(chan claude.Response, 1)

	go claude.StartConversation("foo@example.com", cin, cout)

	for _, q := range []string{"What's the biggest city in Japan?", "Area"} {
		cin <- claude.Request{Prompt: q}
		res := <-cout

		if res.Succeeded {
			log.Infof("ans: %s\n", res.Text)
		} else {
			log.Infof("error: %s\n", res.Text)
		}
	}

	close(cin)

	time.Sleep(time.Second)

	fmt.Println("done")
}
