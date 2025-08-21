package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

	resp, err := http.DefaultClient.Do(req)
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
	httpConfig := &dynamic.HTTPConfiguration{}
	changed := false
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
	if changed {
		configuration := dynamic.Configuration{HTTP: httpConfig}
		return &configuration
	}
	var configuration dynamic.Configuration
	if err := json.Unmarshal(body, &configuration); err != nil {
		return &dynamic.Configuration{}
	}

	return &configuration
}
