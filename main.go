package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world, " + request.Method))
	})

	server := http.Server{
		Handler: mux,
		Addr:    ":3000",
	}
	log.Println("Server start at port 3000")
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server error: %+v\n", err)
		os.Exit(1)
	}
}
