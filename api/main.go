package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	WITH_TICK = true
	NO_TICK   = false
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

func (s server) eventCallbacks(w http.ResponseWriter, r *http.Request) {
	payload := Payload{
		RequestID: "8888-7777-78",
		Data:      []byte("Have a great celebration!"),
	}

	publish("greeting", payload, &s, 1, NO_TICK)
	s.nc.Flush()

	// Terminate connection to NATS
	s.nc.Close()

	err := publish("greeting", payload, &s, 1, NO_TICK)
	if err != nil {
		log.Println(err)
	}

	s.nc.Flush()
}

func (s server) reconnectToServer(w http.ResponseWriter, r *http.Request) {
	_, err := NewServer()
	if err != nil {
		log.Println(err)
	}
}

// func (s server) reconnectTask(w http.ResponseWriter, r *http.Request) {
// 	payload := Payload{
// 		RequestID: "1234-5678-98",
// 		Data:      []byte("try reconnection"),
// 	}

// 	err := publish("greeting", payload, &s, 1, WITH_TICK)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

func (s server) createTask(w http.ResponseWriter, r *http.Request) {
	payload := Payload{
		RequestID: "1234-5678-90",
		Data:      []byte("Can you help me?"),
	}

	publish("greeting", payload, &s, 1, NO_TICK)

	s.nc.Flush()
}

func publish(subject string, payload Payload, s *server, n int, withTick bool) error {

	publishSubject := func(s *server, subject string, payloadJSON []byte) error {
		err := s.nc.Publish(subject, payloadJSON)
		if err != nil {
			return errors.New(fmt.Sprintf("Error on publishing NATS: %v\n", err))
		}

		log.Printf("[PUBLISH] subject '%s' data: \n%s\n", subject, string(payload.Data))

		return nil
	}

	for i := 0; i < n; i++ {
		payload.RequestID = payload.RequestID[:12]
		payload.RequestID = fmt.Sprintf("%s-%d", payload.RequestID, i)
		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return errors.New(fmt.Sprintf("Error on marshalling payload: %v\n", err))

		}
		if withTick {
			for range time.NewTicker(500 * time.Microsecond).C {
				err = publishSubject(s, subject, payloadJSON)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			return publishSubject(s, subject, payloadJSON)
		}
	}

	return nil

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

func NewServer() (server, error) {

	disconnectHandler := func(nc *nats.Conn, err error) {
		if err != nil {
			log.Println("Error while disconnecting:", err)
			return
		}
		log.Println("Disconnected!\n")
	}

	reconnectedHandler := func(nc *nats.Conn) {
		log.Printf("Reconnected to %v!\n", nc.ConnectedUrl())
	}

	closedHandler := func(nc *nats.Conn) {
		log.Printf("Connection closed. Reason %q\n", nc.LastError())
	}

	discoveredServersHandler := func(nc *nats.Conn) {
		log.Printf("Server discovered\n")
	}

	asyncErrorHandler := func(nc *nats.Conn, sub *nats.Subscription, err error) {
		if err != nil {
			log.Printf("Async Error found %q\n.  on subscription %q\n,  %v!\n",
				nc.ConnectedUrl(),
				sub.Subject,
				err,
			)
		}
	}

	var s server
	uri := os.Getenv("NATS_URI")

	opts := nats.DefaultOptions
	// Arbitrarily small reconnecting buffer
	opts.MaxReconnect = 5
	opts.ReconnectBufSize = 256
	opts.Url = uri
	opts.User = "foo"
	opts.Password = "secret"
	opts.Name = "practical-nats-client"
	opts.DisconnectedErrCB = disconnectHandler
	opts.ReconnectedCB = reconnectedHandler
	opts.ClosedCB = closedHandler
	opts.DiscoveredServersCB = discoveredServersHandler
	opts.AsyncErrorCB = asyncErrorHandler

	nc, err := opts.Connect()
	if err != nil {
		return server{}, errors.New(fmt.Sprintf("Error establishing connection to NATS: %v\n", err))
	}

	s.nc = nc
	log.Println("Connected to NATS at:", s.nc.ConnectedUrl())
	return s, nil
}
func main() {

	s, err := NewServer()
	if err != nil {
		log.Println(err)
	}
	http.HandleFunc("/", s.baseRoot)
	http.HandleFunc("/createTask", s.createTask)
	http.HandleFunc("/reconnect", s.reconnectToServer)
	http.HandleFunc("/callbacks", s.eventCallbacks)
	http.HandleFunc("/healthz", s.healthz)

	log.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
