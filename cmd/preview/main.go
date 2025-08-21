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
	"gopkg.in/yaml.v3"
)

func main() {
	cfg := config.CreateConfig()

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
