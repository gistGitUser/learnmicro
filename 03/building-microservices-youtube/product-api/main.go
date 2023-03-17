package main

import (
	"context"
	"github.com/nicholasjackson/env"
	"log"
	"micro03/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"
)

/*
https://www.youtube.com/redirect?event=video_description&redir_token=QUFFLUhqa2dZRlNyb3F0VFdFenRtUXpKRnZzb3hMUDVqZ3xBQ3Jtc0trektjYUZJTjBtYkNmb2dwQlZXQ2Z2dllKbE81RGRJWWxnUTBvNHF0ajVLZTlrMi1rb0JKVUVhaUUyRXNURFpYNUNRVEZ2RnIzMWZzZ3FLRFFRTGJwMTYyN1djSG5vSU5JZTE2R1BwSGpDYWdoeFVENA&q=https%3A%2F%2Fdocs.microsoft.com%2Fen-us%2Fazure%2Farchitecture%2Fbest-practices%2Fapi-design&v=eBeqtmrvVpg
 */

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {

	env.Parse()

	l := log.New(os.Stdout, "products-api ", log.LstdFlags)

	// create the handlers
	ph := handlers.NewProducts(l)

	// create a new serve mux and register the handlers
	sm := http.NewServeMux()
	sm.Handle("/", ph)

	// create a new server
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
