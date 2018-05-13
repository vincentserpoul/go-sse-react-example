package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// SSE hold the references to the different channels
type SSE struct {

	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan []byte

	// New client connections
	newClients chan chan []byte

	// Closed client connections
	closingClients chan chan []byte

	// Client connections registry
	clients map[chan []byte]bool
}

// NewSSE will create a new SSE struct
func NewSSE() (sse *SSE) {
	// Instantiate a broker
	sse = &SSE{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	// Set it running - listening and broadcasting events
	go sse.listen()

	return
}

func (sse *SSE) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// Make sure that the writer supports flushing.
	//
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(chan []byte)

	// Signal the broker that we have a new connection
	sse.newClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		sse.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	notify := rw.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		sse.closingClients <- messageChan
	}()

	for {

		// Write to the ResponseWriter
		// Server Sent Events compatible
		if _, err := fmt.Fprintf(rw, "data: %s\n\n", <-messageChan); err != nil {
			log.Fatalf("%v", err)
		}

		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}

}

func (sse *SSE) listen() {
	for {
		select {
		case s := <-sse.newClients:

			// A new client has connected.
			// Register their message channel
			sse.clients[s] = true
			log.Printf("Client added. %d registered clients", len(sse.clients))
		case s := <-sse.closingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(sse.clients, s)
			log.Printf("Removed client. %d registered clients", len(sse.clients))
		case event := <-sse.Notifier:

			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan := range sse.clients {
				clientMessageChan <- event
			}
		}
	}

}

func main() {

	sse := NewSSE()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/sse", sse)

	go func() {
		for {
			time.Sleep(time.Second * 2)
			eventString := fmt.Sprintf("the time is %v", time.Now())
			log.Println("Receiving event")
			sse.Notifier <- []byte(eventString)
		}
	}()

	if err := http.ListenAndServe("localhost:3000", nil); err != nil {
		log.Fatalf("%v", err)
	}

}
