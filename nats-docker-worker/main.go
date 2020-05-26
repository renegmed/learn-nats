package main

import (
	"encoding/json"
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

		log.Printf("[RECEIVE]\n subject: %s json:\n  request ID: %s\n  data: %v\n",
			m.Subject, payload.RequestID, string(payload.Data))

		if string(payload.Data) == "Can you help me?" {
			reply(m.Reply, "Sure, I would love to help you +++++ WORKER 1", nc)
		}

		//time.Sleep(2000 * time.Millisecond)

	})

	log.Println("----- This worker subscribed to 'greeting' for processing requests...-----")
	log.Println("Server listening on port 8181...")

	http.HandleFunc("/healthz", healthz)
	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal(err)
	}
}

func reply(subject string, message string, nc *nats.Conn) error {
	payload := struct {
		RequestID string `json:"request_id"`
		Data      []byte `json:"data"`
	}{
		RequestID: "2222-3333-99",
		Data:      []byte(message),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Printf("[REPLY] subject '%s' payload: \n%v\n", subject, string(payload.Data))

	err = nc.Publish(subject, payloadJSON)

	return err
}
