//go:build preview
// +build preview

// Package main provides a small preview tool to run the provider and print one configuration snapshot.
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	provider "github.com/zalbiraw/traefikprovider"
	"github.com/zalbiraw/traefikprovider/config"
	"gopkg.in/yaml.v3"
)

func main() {
	cfg := provider.CreateConfig()

	// Override defaults for preview: pull from local Docker-exposed provider1 and provider2
	cfg.PollInterval = "5s"
	cfg.Providers = []config.ProviderConfig{
		{
			Name: "provider1",
			Connection: config.ConnectionConfig{
				Host:    "localhost",
				Port:    8081,
				Path:    "/api/rawdata",
				Timeout: "10s",
			},
			HTTP: &config.HTTPSection{
				Discover: true,
				Routers: &config.RoutersConfig{
					Discover:         true,
					DiscoverPriority: false,
				},
			},
		},
		{
			Name: "provider2",
			Connection: config.ConnectionConfig{
				Host:    "localhost",
				Port:    8082,
				Path:    "/api/rawdata",
				Timeout: "10s",
			},
			HTTP: &config.HTTPSection{
				Discover: true,
				Routers: &config.RoutersConfig{
					Discover:         true,
					DiscoverPriority: false,
				},
			},
		},
	}

	p, err := provider.New(context.Background(), cfg, "preview")
	if err != nil {
		log.Fatalf("init error: %v", err)
	}

	cfgChan := make(chan json.Marshaler, 1)
	if err := p.Provide(cfgChan); err != nil {
		log.Fatalf("provide error: %v", err)
	}

	marshaled := <-cfgChan
	jsonBytes, err := marshaled.MarshalJSON()
	if err != nil {
		log.Fatalf("marshal error: %v", err)
	}
	var asInterface interface{}
	if err := json.Unmarshal(jsonBytes, &asInterface); err != nil {
		log.Fatalf("unmarshal error: %v", err)
	}
	yamlBytes, err := yaml.Marshal(asInterface)
	if err != nil {
		log.Fatalf("yaml marshal error: %v", err)
	}
	os.Stdout.Write(yamlBytes)
	_ = p.Stop()
}
