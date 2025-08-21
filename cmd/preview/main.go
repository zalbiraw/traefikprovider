//go:build preview
// +build preview

// Package main provides a small preview tool to run the provider and print one configuration snapshot.
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	provider "github.com/zalbiraw/traefik-provider"
	"github.com/zalbiraw/traefik-provider/config"
)

func main() {
	cfg := provider.CreateConfig()

	cfg.PollInterval = "5s"

	providerObj := config.ProviderConfig{
		Name: "test",
		Connection: config.ConnectionConfig{
			Host: []string{"dashboard.traefik.localhost"},
			Port: 8080,
		},
	}
	cfg.Providers = append(cfg.Providers, providerObj)

	p, err := provider.New(context.Background(), cfg, "preview")
	if err != nil {
		log.Fatalf("init error: %v", err)
	}

	cfgChan := make(chan json.Marshaler, 1)
	if err := p.Provide(cfgChan); err != nil {
		log.Fatalf("provide error: %v", err)
	}
	// Read a single configuration snapshot and print it.
	marshaled := <-cfgChan
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(marshaled); err != nil {
		log.Fatalf("encode error: %v", err)
	}
	_ = p.Stop()
}
