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
	return parseDynamicConfiguration(body)
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
func parseDynamicConfiguration(body []byte) *dynamic.Configuration {
	var testJson interface{}
	if err := json.Unmarshal(body, &testJson); err != nil {
		return &dynamic.Configuration{}
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return &dynamic.Configuration{}
	}

	var (
		httpConfig *dynamic.HTTPConfiguration
		tcpConfig  *dynamic.TCPConfiguration
		udpConfig  *dynamic.UDPConfiguration
		tlsConfig  *dynamic.TLSConfiguration
		changed    bool
	)

	// HTTP flat keys
	httpConfig = &dynamic.HTTPConfiguration{}
	if routers, ok := raw["routers"]; ok {
		jsonRouters, _ := json.Marshal(routers)
		_ = json.Unmarshal(jsonRouters, &httpConfig.Routers)
		changed = true
	}
	if services, ok := raw["services"]; ok {
		jsonServices, _ := json.Marshal(services)
		_ = json.Unmarshal(jsonServices, &httpConfig.Services)
		changed = true
	}
	if middlewares, ok := raw["middlewares"]; ok {
		jsonMiddlewares, _ := json.Marshal(middlewares)
		_ = json.Unmarshal(jsonMiddlewares, &httpConfig.Middlewares)
		changed = true
	}
	if changed && (len(httpConfig.Routers) > 0 || len(httpConfig.Services) > 0 || len(httpConfig.Middlewares) > 0) {
		return &dynamic.Configuration{HTTP: httpConfig}
	}

	// TCP flat keys
	changed = false
	tcpConfig = &dynamic.TCPConfiguration{}
	if tcpRouters, ok := raw["tcpRouters"]; ok {
		jsonRouters, _ := json.Marshal(tcpRouters)
		_ = json.Unmarshal(jsonRouters, &tcpConfig.Routers)
		changed = true
	}
	if tcpServices, ok := raw["tcpServices"]; ok {
		jsonServices, _ := json.Marshal(tcpServices)
		_ = json.Unmarshal(jsonServices, &tcpConfig.Services)
		changed = true
	}
	if changed && (len(tcpConfig.Routers) > 0 || len(tcpConfig.Services) > 0) {
		return &dynamic.Configuration{TCP: tcpConfig}
	}

	// UDP flat keys
	changed = false
	udpConfig = &dynamic.UDPConfiguration{}
	if udpRouters, ok := raw["udpRouters"]; ok {
		jsonRouters, _ := json.Marshal(udpRouters)
		_ = json.Unmarshal(jsonRouters, &udpConfig.Routers)
		changed = true
	}
	if udpServices, ok := raw["udpServices"]; ok {
		jsonServices, _ := json.Marshal(udpServices)
		_ = json.Unmarshal(jsonServices, &udpConfig.Services)
		changed = true
	}
	if changed && (len(udpConfig.Routers) > 0 || len(udpConfig.Services) > 0) {
		return &dynamic.Configuration{UDP: udpConfig}
	}

	// TLS flat keys
	changed = false
	tlsConfig = &dynamic.TLSConfiguration{}
	if certs, ok := raw["certificates"]; ok {
		jsonCerts, _ := json.Marshal(certs)
		_ = json.Unmarshal(jsonCerts, &tlsConfig.Certificates)
		changed = true
	}
	if options, ok := raw["options"]; ok {
		jsonOpts, _ := json.Marshal(options)
		_ = json.Unmarshal(jsonOpts, &tlsConfig.Options)
		changed = true
	}
	if stores, ok := raw["stores"]; ok {
		jsonStores, _ := json.Marshal(stores)
		_ = json.Unmarshal(jsonStores, &tlsConfig.Stores)
		changed = true
	}
	if changed && (len(tlsConfig.Certificates) > 0 || len(tlsConfig.Options) > 0 || len(tlsConfig.Stores) > 0) {
		return &dynamic.Configuration{TLS: tlsConfig}
	}

	// Nested format (preferred)
	var configuration dynamic.Configuration
	if err := json.Unmarshal(body, &configuration); err != nil {
		return &dynamic.Configuration{}
	}
	return &configuration
}
