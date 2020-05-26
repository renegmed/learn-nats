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

func (s server) reconnectTask(w http.ResponseWriter, r *http.Request) {
	payload := Payload{
		RequestID: "1234-5678-98",
		Data:      []byte("try reconnection"),
	}

	publish("greeting", payload, &s, 1)

}

func (s server) createTask(w http.ResponseWriter, r *http.Request) {
	payload := Payload{
		RequestID: "1234-5678-90",
		Data:      []byte("Can you help me?"),
	}

	request("greeting", payload, &s, 1)

	s.nc.Flush()

	// payload = Payload{
	// 	RequestID: "8888-7777-77",
	// 	Data:      []byte("Have a great celebration!"),
	// }

	// publish(payload, &s, 5)

	// s.nc.Flush()
}

func publish(subject string, payload Payload, s *server, n int) {
	for i := 0; i < n; i++ {
		payload.RequestID = payload.RequestID[:12]
		payload.RequestID = fmt.Sprintf("%s-%d", payload.RequestID, i)
		payloadJSON, err := json.Marshal(payload)

		for range time.NewTicker(500 * time.Microsecond).C {
			if s.nc.IsClosed() {
				log.Fatal("Disconnected forever! Exiting....")
			}
			if s.nc.IsReconnecting() {
				log.Println("Disconnected temporarily, skipping for now")
				continue
			}

			err = s.nc.Publish(subject, payloadJSON)
			if err != nil {
				log.Fatalf("Error on making NATS publish %d: %v\n", i, err)
			}

			log.Printf("[PUBLISH] subject '%s' %d. data: \n%s\n", subject, i, string(payload.Data))
		}
	}

}

func request(subject string, payload Payload, s *server, n int) {
	for i := 0; i < n; i++ {
		payload.RequestID = payload.RequestID[:12]
		payload.RequestID = fmt.Sprintf("%s-%d", payload.RequestID, i)
		payloadJSON, err := json.Marshal(payload)

		log.Printf("[REQUEST] subject '%s' %d. data: \n%v\n", subject, i, string(payload.Data))

		response, err := s.nc.Request(subject, payloadJSON, 1*time.Second)
		if err != nil {
			log.Printf("Error while making NATS request %d: %v\n", i, err)
		}
		p, err := processResponse(response)
		if err != nil {
			log.Println("Error on unmarshal of response", err)
		}
		log.Println("[RESPONSE] ", string(p.Data))
	}
}

func processResponse(msg *nats.Msg) (*Payload, error) {
	payload := &Payload{}
	err := json.Unmarshal([]byte(msg.Data), payload)
	return payload, err
}
func (s server) healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func main() {
	var s server
	var err error
	uri := os.Getenv("NATS_URI")

	opts := nats.DefaultOptions
	// Arbitrarily small reconnecting buffer
	opts.ReconnectBufSize = 256
	opts.Url = uri
	opts.User = "foo"
	opts.Password = "secret"
	opts.Name = "practical-nats-client"
	nc, err := opts.Connect()

	// nc, err := nats.Connect(uri, nats.Name("practical-nats-client"),
	// 	nats.UserInfo("foo", "secret"),
	// )

	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}

	s.nc = nc

	log.Println("Connected to NATS at:", s.nc.ConnectedUrl())
	http.HandleFunc("/", s.baseRoot)
	http.HandleFunc("/createTask", s.createTask)
	http.HandleFunc("/reconnectTask", s.reconnectTask)
	http.HandleFunc("/healthz", s.healthz)

	log.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
