package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

type server struct {
	nc *nats.Conn
}

func (s server) baseRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Basic NATS based microservice example v0.0.1")
}

func (s server) createTask(w http.ResponseWriter, r *http.Request) {
	err := s.nc.Publish("greeting", []byte("hello world"))
	if err != nil {
		log.Println("Error making NATS request:", err)
	}
	log.Print("+++ Published 'hello  world' to subject 'greeting'")
}

func (s server) healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func main() {
	var s server
	var err error
	uri := os.Getenv("NATS_URI")

	nc, err := nats.Connect(uri, nats.Name("practical-nats-client"),
		nats.UserInfo("foo", "secret"),
	)

	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}

	s.nc = nc

	log.Println("Connected to NATS at:", s.nc.ConnectedUrl())
	http.HandleFunc("/", s.baseRoot)
	http.HandleFunc("/createTask", s.createTask)
	http.HandleFunc("/healthz", s.healthz)

	log.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
