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
	"github.com/zalbiraw/traefik-provider/internal"
	"github.com/zalbiraw/traefik-provider/internal/httpclient"
)

type Provider struct {
	name         string
	pollInterval time.Duration
	config       *config.Config
	cancel       func()
}

func New(ctx context.Context, config *config.Config, name string) (*Provider, error) {
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

	if len(config.Providers) == 0 {
		return nil, fmt.Errorf("at least one ProviderConfig is required")
	}
	for i, p := range config.Providers {
		if p.Name == "" {
			return nil, fmt.Errorf("provider[%d]: Name is required", i)
		}
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

func (p *Provider) Init() error {
	if p.pollInterval <= 0 {
		return fmt.Errorf("poll interval must be greater than 0")
	}

	return nil
}

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

	for {
		select {
		case <-ticker.C:
			var configs []*dynamic.Configuration
			for i := range p.config.Providers {
				cfg := httpclient.GenerateConfiguration(&p.config.Providers[i])
				configs = append(configs, cfg)
			}
			merged := internal.MergeConfigurations(configs...)
			cfgChan <- &dynamic.JSONPayload{Configuration: merged}
		case <-ctx.Done():
			return
		}
	}
}

func (p *Provider) Stop() error {
	p.cancel()
	return nil
}
