package fontytelemetry

import (
	"encoding/json"
	"log"
)

// Configuration is a struct describing the application's configuration.
type Configuration struct {
	ProjectID          string
	PubsubTopic        string
	PubsubSubscription string
}

// Config is a stateful variable containing the loaded application configuration.
var Config Configuration

// LoadConfig loads the application configuration.
func LoadConfig() *Configuration {
	data, err := Asset("cfg/config.json")

	// Decode JSON body
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Failure to load configuration: %v", err)
	}
	return &Config
}
