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

type Payload struct {
	RequestID string `json:"request_id"`
	Data      []byte `json:"data"`
}

func (s server) baseRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Basic NATS based microservice example v0.0.1")
}

func (s server) createTask(w http.ResponseWriter, r *http.Request) {
	payload := Payload{
		RequestID: "1234-5678-90",
		Data:      []byte("Happy birthday to you!"),
	}

	publish(payload, &s, 10)

	s.nc.Flush()

	payload = Payload{
		RequestID: "8888-7777-77",
		Data:      []byte("Have a great celebration!"),
	}

	publish(payload, &s, 5)

	s.nc.Flush()
}

func publish(payload Payload, s *server, n int) {
	for i := 0; i < n; i++ {
		payload.RequestID = payload.RequestID[:12]
		payload.RequestID = fmt.Sprintf("%s-%d", payload.RequestID, i)
		payloadJSON, err := json.Marshal(payload)

		err = s.nc.Publish("greeting", payloadJSON)
		if err != nil {
			log.Printf("Error on making NATS request %d: %v\n", i, err)
		}
		log.Print("[Published] subject 'greeting' %d. data: \n%s\n", i, string(payloadJSON))
	}
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
