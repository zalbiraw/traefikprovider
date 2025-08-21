// Package traefik-provider contains a demo of the provider's plugin.
package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/traefik/genconf/dynamic"

	"github.com/zalbiraw/traefik-provider/config"
	"github.com/zalbiraw/traefik-provider/provider/httpclient"
)


// Provider a simple provider plugin.
type Provider struct {
	name         string
	pollInterval time.Duration
	config       *config.Config
	cancel       func()
}

// New creates a new Provider plugin.
func New(ctx context.Context, config *config.Config, name string) (*Provider, error) {
	// Validate PollInterval
	if config.PollInterval == "" {
		return nil, fmt.Errorf("PollInterval is required")
	}
	pi, err := time.ParseDuration(config.PollInterval)
	if err != nil {
		return nil, fmt.Errorf("invalid PollInterval: %w", err)
	}
	if pi <= 0 {
		return nil, fmt.Errorf("PollInterval must be greater than 0")
	}

	// Validate Providers
	if len(config.Providers) == 0 {
		return nil, fmt.Errorf("at least one ProviderConfig is required")
	}
	for i, p := range config.Providers {
		if p.Name == "" {
			return nil, fmt.Errorf("provider[%d]: Name is required", i)
		}
		// Example: validate Connection fields
		if len(p.Connection.Host) == 0 {
			return nil, fmt.Errorf("provider[%d]: Connection.Host is required", i)
		}
		if p.Connection.Port == 0 {
			return nil, fmt.Errorf("provider[%d]: Connection.Port is required", i)
		}
	}

	return &Provider{
		name:         name,
		pollInterval: pi,
		config:       config,
	}, nil
}

// Init the provider.
func (p *Provider) Init() error {
	if p.pollInterval <= 0 {
		return fmt.Errorf("poll interval must be greater than 0")
	}

	return nil
}

// Provide creates and send dynamic configuration.
func (p *Provider) Provide(cfgChan chan<- json.Marshaler) error {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Print(err)
			}
		}()

		p.loadConfiguration(ctx, cfgChan)
	}()

	return nil
}

func (p *Provider) loadConfiguration(ctx context.Context, cfgChan chan<- json.Marshaler) {
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	// TODO: Support multiple providers if needed. For now, use the first one.
	providerCfg := &config.ProviderConfig{}
	if len(p.config.Providers) > 0 {
		providerCfg = &p.config.Providers[0]
	}

	for {
		select {
		case <-ticker.C:
			configuration := httpclient.GenerateConfiguration(providerCfg)
			cfgChan <- &dynamic.JSONPayload{Configuration: configuration}
		case <-ctx.Done():
			return
		}
	}
}

// Stop to stop the provider and the related go routines.
func (p *Provider) Stop() error {
	p.cancel()
	return nil
}
