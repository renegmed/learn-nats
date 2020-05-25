package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func main() {
	uri := os.Getenv("NATS_URI")
	var err error
	var nc *nats.Conn

	nc, err = nats.Connect(uri, nats.Name("practical-nats-client"))
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}

	log.Println("Connected to NATS at:", nc.ConnectedUrl())

	nc.Subscribe("greeting", func(m *nats.Msg) {
		log.Printf("Got subject '%s' message: \n\t %s\n", m.Subject, string(m.Data))
	})

	log.Println("Worker subscribed to 'greeting' for processing requests...")
	log.Println("Server listening on port 8181...")

	http.HandleFunc("/healthz", healthz)
	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal(err)
	}
}
