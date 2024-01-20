package main

import (
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

func sse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	log.Println("client connected")

	clientChan := make(chan string)
	id := uuid.New().String()
	Clients[id] = Client{channel: clientChan}

	// Clean up client when connection closed
	defer func() {
		delete(Clients, id)
		close(clientChan)
	}()

	closeNotify := r.Context().Done()
	w.WriteHeader(http.StatusOK)
	for {
		select {
		case <-closeNotify:
			log.Println("Client closed")
			return
		case data := <-clientChan:
			log.Println(data)
			w.Write([]byte(data))
			w.(http.Flusher).Flush()
		}
	}

}

func broadcast(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad data"))
	}

	// Send data to all clients
	func(data string) {
		for _, v := range Clients {
			v.channel <- data
		}
	}(string(body))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("message sent"))
}
