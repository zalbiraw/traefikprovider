package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/traefik/genconf/dynamic"

	"github.com/zalbiraw/traefik-provider/config"
)

// GenerateConfiguration fetches and parses the dynamic configuration from the remote provider.
func GenerateConfiguration(providerCfg *config.ProviderConfig) *dynamic.Configuration {
	if providerCfg.Connection.Host == "" || providerCfg.Connection.Port == 0 || providerCfg.Connection.Path == "" {
		return &dynamic.Configuration{}
	}

	url := buildProviderURL(providerCfg)
	req := buildProviderRequest(url, providerCfg.Connection.Headers)

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
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
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
	return fmt.Sprintf("http://%s:%d%s", host, port, path)
}

// buildProviderRequest creates an HTTP GET request with headers.
func buildProviderRequest(url string, headers map[string]string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if hostHeader, ok := headers["Host"]; ok {
		req.Host = hostHeader
	}
	return req
}

// parseDynamicConfiguration parses the response body into a dynamic.Configuration struct.
func parseDynamicConfiguration(body []byte, providerCfg *config.ProviderConfig) (*dynamic.Configuration, error) {
	var testJson interface{}
	if err := json.Unmarshal(body, &testJson); err != nil {
		return &dynamic.Configuration{}, fmt.Errorf("error unmarshaling response body to testJson: %w", err)
	}
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

	// HTTP
	if providerCfg.HTTP.Discover {
		parseHTTPConfig(raw, httpConfig, providerCfg.HTTP, providerCfg.Filter, providerCfg.Tunnels)
	}

	// TCP
	if providerCfg.TCP.Discover {
		parseTCPConfig(raw, tcpConfig, providerCfg.TCP, providerCfg.Filter, providerCfg.Tunnels)
	}

	// UDP
	if providerCfg.UDP.Discover {
		parseUDPConfig(raw, udpConfig, providerCfg.UDP, providerCfg.Filter, providerCfg.Tunnels)
	}

	// TLS
	if providerCfg.TLS.Discover {
		parseTLSConfig(raw, tlsConfig, providerCfg.TLS)
	}

	cfg := &dynamic.Configuration{
		HTTP: httpConfig,
		TCP:  tcpConfig,
		UDP:  udpConfig,
		TLS:  tlsConfig,
	}
	return cfg, nil
}
