package main

import (
	"fmt"
	"net"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	server := &http.Server{Addr: ":80"}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}

	fmt.Println("Listening on port 80")

	err = server.Serve(listener)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
