package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func main() {
	uri := os.Getenv("NATS_URI")
	var err error
	var nc *nats.Conn

	nc, err = nats.Connect(uri, nats.Name("practical-nats-worker"),
		nats.UserInfo("foo", "secret"),
	)

	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}

	log.Println("Connected to NATS at:", nc.ConnectedUrl())

	// server balances requests randomly among the members of the group, workers-group
	nc.QueueSubscribe("greeting", "workers-group", func(m *nats.Msg) {

		payload := struct {
			RequestID string `json:"request_id"`
			Data      []byte `json:"data"`
		}{}

		err := json.Unmarshal([]byte(m.Data), &payload)
		if err != nil {
			log.Fatalf("Error on unmarshalling payload: %v", err)
		}

		log.Printf("[Received]\n subject: %s json:\n  request ID: %s\n  data: %v\n",
			m.Subject, payload.RequestID, string(payload.Data))

		time.Sleep(500 * time.Millisecond)

	})

	log.Println("----- This worker subscribed to 'greeting' for processing requests...-----")
	log.Println("Server listening on port 8181...")

	http.HandleFunc("/healthz", healthz)
	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal(err)
	}
}
