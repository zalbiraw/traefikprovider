package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/traefik/genconf/dynamic"

	"github.com/zalbiraw/traefik-provider/config"
)

// generateConfiguration fetches and parses the dynamic configuration from the remote provider.
func GenerateConfiguration(providerCfg *config.ProviderConfig) *dynamic.Configuration {
	if providerCfg.Connection.Host == "" {
		log.Print("No host configured for provider")
		return &dynamic.Configuration{}
	}

	host := providerCfg.Connection.Host
	port := providerCfg.Connection.Port
	path := providerCfg.Connection.Path

	url := fmt.Sprintf("http://%s:%d%s", host, port, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request for %s: %v", url, err)
		return &dynamic.Configuration{}
	}
	// Add headers from config
	log.Printf("Sending headers:")
	for k, v := range providerCfg.Connection.Headers {
		log.Printf("  %s: %s", k, v)
		req.Header.Set(k, v)
	}
	// Set req.Host if 'Host' header is present
	if hostHeader, ok := providerCfg.Connection.Headers["Host"]; ok {
		log.Printf("Setting req.Host to: %s", hostHeader)
		req.Host = hostHeader
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to fetch config from %s: %v", url, err)
		return &dynamic.Configuration{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return &dynamic.Configuration{}
	}
	log.Printf("Raw response body: %s", string(body))

	var testJson interface{}
	if err := json.Unmarshal(body, &testJson); err != nil {
		log.Printf("Response is not valid JSON.")
		return &dynamic.Configuration{}
	}

	var configuration dynamic.Configuration
	if err := json.Unmarshal(body, &configuration); err != nil {
		log.Printf("Failed to parse config JSON: %v", err)
		return &dynamic.Configuration{}
	}

	return &configuration
}
