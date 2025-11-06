// Package httpclient fetches provider data over HTTP and converts it into Traefik dynamic configuration.
package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
	"github.com/zalbiraw/traefikprovider/internal/parsers"
)

// GenerateConfiguration fetches and parses the dynamic configuration from the remote provider.
func GenerateConfiguration(providerCfg *config.ProviderConfig) *dynamic.Configuration {
	if providerCfg.Connection.Host == "" || providerCfg.Connection.Port == 0 || providerCfg.Connection.Path == "" {
		return &dynamic.Configuration{}
	}

	url := buildProviderURL(providerCfg)
	req := buildProviderRequest(url, providerCfg.Connection.Host, providerCfg.Connection.Headers)

	client := http.DefaultClient
	if providerCfg.Connection.Timeout != "" {
		if d, err := time.ParseDuration(providerCfg.Connection.Timeout); err == nil {
			client = &http.Client{Timeout: d}
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return &dynamic.Configuration{}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return &dynamic.Configuration{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &dynamic.Configuration{}
	}
	cfg, err := parseDynamicConfiguration(body, providerCfg)
	if err != nil {
		return cfg
	}

	return cfg
}

// buildProviderURL constructs the URL for the provider endpoint.
func buildProviderURL(cfg *config.ProviderConfig) string {
	host := cfg.Connection.Host
	port := cfg.Connection.Port
	path := cfg.Connection.Path
	scheme := "http"
	if cfg.Connection.HTTPS {
		scheme = "https"
	}
	hostPort := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	return scheme + "://" + hostPort + path
}

// buildProviderRequest creates an HTTP GET request with headers and host.
func buildProviderRequest(url, host string, headers map[string]string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Host = host
	return req
}

// parseDynamicConfiguration parses the response body into a dynamic.Configuration struct.
func parseDynamicConfiguration(body []byte, providerCfg *config.ProviderConfig) (*dynamic.Configuration, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return &dynamic.Configuration{}, fmt.Errorf("error unmarshaling response body to raw map: %w", err)
	}

	var (
		httpConfig = &dynamic.HTTPConfiguration{}
		tcpConfig  = &dynamic.TCPConfiguration{}
		udpConfig  = &dynamic.UDPConfiguration{}
		tlsConfig  = &dynamic.TLSConfiguration{}
	)

	ensureProviderDefaults(providerCfg)

	// HTTP
	if providerCfg.HTTP.Discover {
		parsers.ParseHTTPConfig(raw, httpConfig, providerCfg.HTTP, providerCfg.Matcher, providerCfg.Tunnels)
	}

	// TCP
	if providerCfg.TCP.Discover {
		parsers.ParseTCPConfig(raw, tcpConfig, providerCfg.TCP, providerCfg.Matcher, providerCfg.Tunnels)
	}

	// UDP
	if providerCfg.UDP.Discover {
		parsers.ParseUDPConfig(raw, udpConfig, providerCfg.UDP, providerCfg.Matcher)
	}

	// TLS
	if providerCfg.TLS.Discover {
		parsers.ParseTLSConfig(raw, tlsConfig, providerCfg.TLS)
	}

	cfg := &dynamic.Configuration{
		HTTP: httpConfig,
		TCP:  tcpConfig,
		UDP:  udpConfig,
		TLS:  tlsConfig,
	}

	return cfg, nil
}

// ensureProviderDefaults ensures non-nil sections with default Discover values.
func ensureProviderDefaults(providerCfg *config.ProviderConfig) {
	if providerCfg.HTTP == nil {
		providerCfg.HTTP = &config.HTTPSection{Discover: true}
	}
	if providerCfg.TCP == nil {
		providerCfg.TCP = &config.TCPSection{Discover: true}
	}
	if providerCfg.UDP == nil {
		providerCfg.UDP = &config.UDPSection{Discover: true}
	}
	if providerCfg.TLS == nil {
		providerCfg.TLS = &config.TLSSection{Discover: true}
	}
	if providerCfg.Tunnels == nil {
		providerCfg.Tunnels = []config.TunnelConfig{}
	}
}
