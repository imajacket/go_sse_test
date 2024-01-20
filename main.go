package main

import (
	"log"
	"net/http"
)

type Client struct {
	channel chan string
}

var Clients map[string]Client

func main() {
	Clients = make(map[string]Client)

	mux := http.NewServeMux()

	mux.HandleFunc("/sse", sse)
	mux.HandleFunc("/broadcast", broadcast)

	log.Fatal(http.ListenAndServe(":3030", mux))
}
