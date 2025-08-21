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
	fmt.Printf("[DEBUG] Provider config: %+v\n", providerCfg)
	if len(providerCfg.Connection.Host) == 0 {
		return &dynamic.Configuration{}
	}

	url := buildProviderURL(providerCfg)
	req := buildProviderRequest(url, providerCfg.Connection.Headers)
	if req == nil {
		return &dynamic.Configuration{}
	}

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &dynamic.Configuration{}
	}
	cfg, err := parseDynamicConfiguration(body, providerCfg)
	if err != nil {
		fmt.Printf("[DEBUG] Error in parseDynamicConfiguration: %v\n", err)
		return cfg
	}
	fmt.Printf("[DEBUG] Final configuration from GenerateConfiguration: %+v\n", cfg)
	return cfg
}

// buildProviderURL constructs the URL for the provider endpoint.
func buildProviderURL(cfg *config.ProviderConfig) string {
	host := cfg.Connection.Host[0]
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
	fmt.Printf("[DEBUG] Received body: %s\n", string(body))
	var testJson interface{}
	if err := json.Unmarshal(body, &testJson); err != nil {
		return &dynamic.Configuration{}, fmt.Errorf("error unmarshaling response body to testJson: %w", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return &dynamic.Configuration{}, fmt.Errorf("error unmarshaling response body to raw map: %w", err)
	}
	fmt.Printf("[DEBUG] Parsed raw map: %+v\n", raw)

	var (
		httpConfig *dynamic.HTTPConfiguration
		tcpConfig  *dynamic.TCPConfiguration
		udpConfig  *dynamic.UDPConfiguration
		tlsConfig  *dynamic.TLSConfiguration
	)

	// HTTP
	httpConfig = &dynamic.HTTPConfiguration{}
	if providerCfg.HTTP != nil && providerCfg.HTTP.Discover {
		if err := parseHTTPConfig(raw, httpConfig); err != nil {
			return &dynamic.Configuration{}, err
		}
	}

	// TCP
	tcpConfig = &dynamic.TCPConfiguration{}
	if providerCfg.TCP != nil && providerCfg.TCP.Discover {
		if err := parseTCPConfig(raw, tcpConfig); err != nil {
			return &dynamic.Configuration{}, err
		}
	}

	// UDP
	udpConfig = &dynamic.UDPConfiguration{}
	if providerCfg.UDP != nil && providerCfg.UDP.Discover {
		if err := parseUDPConfig(raw, udpConfig); err != nil {
			return &dynamic.Configuration{}, err
		}
	}

	// TLS
	tlsConfig = &dynamic.TLSConfiguration{}
	if providerCfg.TLS != nil && providerCfg.TLS.Discover {
		if err := parseTLSConfig(raw, tlsConfig); err != nil {
			return &dynamic.Configuration{}, err
		}
	}

	cfg := &dynamic.Configuration{
		HTTP: httpConfig,
		TCP:  tcpConfig,
		UDP:  udpConfig,
		TLS:  tlsConfig,
	}
	fmt.Printf("[DEBUG] Returning configuration: %+v\n", cfg)
	return cfg, nil
}
