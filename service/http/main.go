package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vincentserpoul/go-sse-react-example/pkg/sse"
)

func main() {

	s := sse.NewSSE()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/sse", s)

	go func() {
		for {
			time.Sleep(time.Second * 2)
			eventString := fmt.Sprintf("the time is %v", time.Now())
			log.Println("Receiving event")
			s.Notifier <- []byte(eventString)
		}
	}()

	if err := http.ListenAndServe("localhost:3000", nil); err != nil {
		log.Fatalf("%v", err)
	}

}
