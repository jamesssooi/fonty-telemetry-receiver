package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/jamesssooi/fonty-telemetry-receiver/pkg/fontytelemetry"
	"github.com/jamesssooi/fonty-telemetry-receiver/pkg/iso8601"

	"cloud.google.com/go/pubsub"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

var ctx = context.Background()

func main() {
	// Load configuration
	_ = fontytelemetry.LoadConfig()

	// Make router
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1", endpointV1)

	// Run server
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", "localhost", fontytelemetry.Config.Port), router))
}

// TelemetryData represents an incoming telemetry event from a fonty client.
type TelemetryData struct {
	IPAddress       string          `json:"ip_address"`
	ServerTimestamp int64           `json:"server_timestamp"`
	Timestamp       iso8601.ISO8601 `json:"timestamp"`
	StatusCode      int             `json:"status_code"`
	EventType       string          `json:"event_type"`
	ExecutionTime   float32         `json:"execution_time"`
	FontyVersion    string          `json:"fonty_version"`
	OSFamily        string          `json:"os_family"`
	OSVersion       string          `json:"os_version"`
	PythonVersion   string          `json:"python_version"`
	Data            interface{}     `json:"data"`
}

func endpointV1(w http.ResponseWriter, r *http.Request) {
	var data TelemetryData

	// Decode JSON body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Fatalf("Failed to decode JSON data: %v", err)
	}

	// Add server timetamp
	data.ServerTimestamp = time.Now().Unix()

	// Add IP address. Needed to process geographic area later, will be removed in
	// the ETL pipeline before storage.
	ipaddr := r.Header.Get("X-Real-IP")
	if ipaddr == "" {
		ipaddr, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	data.IPAddress = ipaddr

	// Recode JSON body
	o, _ := json.Marshal(data)

	// Publish to Google PubSub
	client, ctx := fontytelemetry.GetClient()
	topic := client.Topic(fontytelemetry.Config.PubsubTopic)
	res := topic.Publish(ctx, &pubsub.Message{Data: []byte(o)})

	// Await response
	id, err := res.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	// Send response to client
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"status\": \"ok\", \"id\": \"%v\"}", id)

	// Log to stdout
	log.Printf("%v - %v/%v", data.EventType, fontytelemetry.Config.PubsubTopic, id)
}
