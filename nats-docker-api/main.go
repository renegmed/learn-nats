package main

import (
	"encoding/json"
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
	payload := struct {
		RequestID string `json:"request_id"`
		Data      []byte `json:"data"`
	}{
		RequestID: "1234-5678-90",
		Data:      []byte("Happy birthday to you!"),
	}

	payloadJSON, err := json.Marshal(payload)

	//log.Printf("Json data to published:\n %s\n", string(payloadJSON))

	err = s.nc.Publish("greeting", payloadJSON)
	if err != nil {
		log.Println("Error on making NATS request:", err)
	}
	log.Print("[Published] subject 'greeting' data: \n%s\n", string(payloadJSON))
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
