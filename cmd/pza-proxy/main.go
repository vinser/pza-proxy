package main

import (
	"log"
	"net/http"

	"github.com/vinser/pza-proxy/internal/api"
	"github.com/vinser/pza-proxy/internal/config"
)

func main() {
	if err := config.Load("config.yaml"); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("GLOBAL LOG: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	})

	chat := api.NewChatHandler()

	for _, path := range config.C.Routing.ChatPaths {
		mux.Handle(path, chat)
	}

	server := &http.Server{
		Addr:    config.C.Server.Listen,
		Handler: mux,
	}

	log.Printf("pza-proxy listening on %s", config.C.Server.Listen)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
