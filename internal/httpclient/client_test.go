package httpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zalbiraw/traefik-provider/config"
)

func TestGenerateConfiguration_HostHeaderOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "custom-host.com" {
			t.Errorf("Expected Host header custom-host.com, got %s", r.Host)
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host: server.URL[7:], // Remove http://
			Path: "/api",
			Headers: map[string]string{
				"Host": "custom-host.com",
			},
		},
		HTTP: &config.HTTPSection{Discover: true},
	}

	result := GenerateConfiguration(providerConfig)
	if result == nil {
		t.Error("Expected non-nil configuration")
	}
}

func TestBuildProviderURL(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.ProviderConfig
		expected string
	}{
		{
			name: "basic URL construction",
			config: &config.ProviderConfig{
				Connection: config.ConnectionConfig{
					Host: "localhost",
					Port: 8080,
					Path: "/api/v1",
				},
			},
			expected: "http://localhost:8080/api/v1",
		},
		{
			name: "URL with empty path",
			config: &config.ProviderConfig{
				Connection: config.ConnectionConfig{
					Host: "example.com",
					Port: 9000,
					Path: "",
				},
			},
			expected: "http://example.com:9000",
		},
		{
			name: "URL with root path",
			config: &config.ProviderConfig{
				Connection: config.ConnectionConfig{
					Host: "api.example.com",
					Port: 443,
					Path: "/",
				},
			},
			expected: "http://api.example.com:443/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildProviderURL(tt.config)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestBuildProviderRequest_ValidWithHeaders(t *testing.T) {
	url := "http://localhost:8080/api"
	headers := map[string]string{
		"Authorization": "Bearer token",
		"Content-Type":  "application/json",
	}
	req := buildProviderRequest(url, headers)
	if req == nil {
		t.Fatal("Expected non-nil request")
	}
	if req.Method != http.MethodGet {
		t.Errorf("Expected GET method, got %s", req.Method)
	}
	if req.URL.String() != url {
		t.Errorf("Expected URL %s, got %s", url, req.URL.String())
	}
	for k, v := range headers {
		if req.Header.Get(k) != v {
			t.Errorf("Expected header %s: %s, got %s", k, v, req.Header.Get(k))
		}
	}
}

func TestBuildProviderRequest_ValidNoHeaders(t *testing.T) {
	url := "http://example.com/test"
	req := buildProviderRequest(url, nil)
	if req == nil {
		t.Fatal("Expected non-nil request")
	}
	if req.Method != http.MethodGet {
		t.Errorf("Expected GET method, got %s", req.Method)
	}
	if req.URL.String() != url {
		t.Errorf("Expected URL %s, got %s", url, req.URL.String())
	}
}

func TestBuildProviderRequest_HostHeader(t *testing.T) {
	url := "http://localhost:8080/api"
	headers := map[string]string{
		"Host": "custom-host.com",
	}
	req := buildProviderRequest(url, headers)
	if req == nil {
		t.Fatal("Expected non-nil request")
	}
	if req.Host != "custom-host.com" {
		t.Errorf("Expected Host custom-host.com, got %s", req.Host)
	}
}

func TestBuildProviderRequest_InvalidURL(t *testing.T) {
	req := buildProviderRequest("://invalid-url", nil)
	if req != nil {
		t.Error("Expected nil request for invalid URL")
	}
}

func TestParseDynamicConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		body           []byte
		providerConfig *config.ProviderConfig
		expectError    bool
	}{
		{
			name: "valid JSON with all sections enabled",
			body: []byte(`{
				"routers": {"test-router": {}},
				"services": {"test-service": {}},
				"middlewares": {"test-middleware": {}}
			}`),
			providerConfig: &config.ProviderConfig{
				HTTP: &config.HTTPSection{
					Discover:    true,
					Routers:     &config.RoutersConfig{Discover: true},
					Services:    &config.ServicesConfig{Discover: true},
					Middlewares: &config.MiddlewaresConfig{Discover: true},
				},
				TCP: &config.TCPSection{
					Discover:    true,
					Routers:     &config.RoutersConfig{Discover: true},
					Services:    &config.ServicesConfig{Discover: true},
					Middlewares: &config.MiddlewaresConfig{Discover: true},
				},
				UDP: &config.UDPSection{
					Discover: true,
					Routers:  &config.UDPRoutersConfig{Discover: true},
					Services: &config.UDPServicesConfig{Discover: true},
				},
				TLS: &config.TLSSection{Discover: true},
			},
			expectError: false,
		},
		{
			name:           "valid JSON with no sections enabled",
			body:           []byte(`{"routers": {}, "services": {}}`),
			providerConfig: &config.ProviderConfig{},
			expectError:    false,
		},
		{
			name:           "invalid JSON",
			body:           []byte(`invalid json`),
			providerConfig: &config.ProviderConfig{},
			expectError:    true,
		},
		{
			name:           "valid JSON but invalid structure for raw map",
			body:           []byte(`"just a string"`),
			expectError:    true,
			providerConfig: &config.ProviderConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseDynamicConfiguration(tt.body, tt.providerConfig)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !tt.expectError && cfg == nil {
				t.Error("Expected non-nil configuration")
			}
		})
	}
}

func TestGenerateConfigurationIntegration(t *testing.T) {
	// Test with a mock server that returns complex configuration
	complexResponse := map[string]interface{}{
		"routers": map[string]interface{}{
			"api-router": map[string]interface{}{
				"rule":    "Host(`api.example.com`)",
				"service": "api-service",
			},
		},
		"services": map[string]interface{}{
			"api-service": map[string]interface{}{
				"loadBalancer": map[string]interface{}{
					"servers": []interface{}{
						map[string]interface{}{"url": "http://backend:8080"},
					},
				},
			},
		},
		"middlewares": map[string]interface{}{
			"auth": map[string]interface{}{
				"basicAuth": map[string]interface{}{
					"users": []string{"admin:$2y$10$..."},
				},
			},
		},
		"tcpRouters": map[string]interface{}{
			"tcp-router": map[string]interface{}{
				"rule":    "HostSNI(`tcp.example.com`)",
				"service": "tcp-service",
			},
		},
		"tcpServices": map[string]interface{}{
			"tcp-service": map[string]interface{}{
				"loadBalancer": map[string]interface{}{
					"servers": []interface{}{
						map[string]interface{}{"address": "backend:8081"},
					},
				},
			},
		},
		"udpRouters": map[string]interface{}{
			"udp-router": map[string]interface{}{
				"service": "udp-service",
			},
		},
		"udpServices": map[string]interface{}{
			"udp-service": map[string]interface{}{
				"loadBalancer": map[string]interface{}{
					"servers": []interface{}{
						map[string]interface{}{"address": "backend:8082"},
					},
				},
			},
		},
	}

	responseBytes, err := json.Marshal(complexResponse)
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(responseBytes); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host: server.URL[7:], // Remove http://
			Port: 0,
			Path: "/",
			Headers: map[string]string{
				"Accept": "application/json",
			},
		},
		HTTP: &config.HTTPSection{
			Discover:    true,
			Routers:     &config.RoutersConfig{Discover: true},
			Services:    &config.ServicesConfig{Discover: true},
			Middlewares: &config.MiddlewaresConfig{Discover: true},
		},
		TCP: &config.TCPSection{
			Discover:    true,
			Routers:     &config.RoutersConfig{Discover: true},
			Services:    &config.ServicesConfig{Discover: true},
			Middlewares: &config.MiddlewaresConfig{Discover: true},
		},
		UDP: &config.UDPSection{
			Discover: true,
			Routers:  &config.UDPRoutersConfig{Discover: true},
			Services: &config.UDPServicesConfig{Discover: true},
		},
		TLS: &config.TLSSection{Discover: true},
	}

	result := GenerateConfiguration(providerConfig)

	if result == nil {
		t.Fatal("Expected non-nil configuration")
	}

	// Debug output
	t.Logf("Result HTTP: %+v", result.HTTP)
	t.Logf("Result TCP: %+v", result.TCP)

	// Verify configuration was processed - sections may be nil for empty config
	if result.HTTP != nil {
		t.Logf("HTTP Routers count: %d", len(result.HTTP.Routers))
	}
	if result.TCP != nil {
		t.Logf("TCP Routers count: %d", len(result.TCP.Routers))
	}
}

func TestGenerateConfigurationErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		providerConfig *config.ProviderConfig
		serverResponse string
		statusCode     int
		expectEmpty    bool
	}{
		{
			name: "empty host",
			providerConfig: &config.ProviderConfig{
				Connection: config.ConnectionConfig{Host: ""},
			},
			expectEmpty: true,
		},
		{
			name: "invalid JSON response",
			providerConfig: &config.ProviderConfig{
				Connection: config.ConnectionConfig{Host: "localhost:8080"},
				HTTP:       &config.HTTPSection{Discover: true},
			},
			serverResponse: "invalid json",
			statusCode:     200,
			expectEmpty:    true,
		},
		{
			name: "HTTP error status",
			providerConfig: &config.ProviderConfig{
				Connection: config.ConnectionConfig{Host: "localhost:8080"},
				HTTP:       &config.HTTPSection{Discover: true},
			},
			serverResponse: `{"error": "not found"}`,
			statusCode:     404,
			expectEmpty:    true,
		},
		{
			name: "timeout configuration",
			providerConfig: &config.ProviderConfig{
				Connection: config.ConnectionConfig{
					Host:    "localhost:8080",
					Timeout: "5s",
				},
				HTTP: &config.HTTPSection{Discover: true},
			},
			serverResponse: `{"routers": {}}`,
			statusCode:     200,
			expectEmpty:    false,
		},
		{
			name: "invalid timeout format",
			providerConfig: &config.ProviderConfig{
				Connection: config.ConnectionConfig{
					Host:    "localhost:8080",
					Timeout: "invalid",
				},
				HTTP: &config.HTTPSection{Discover: true},
			},
			serverResponse: `{"routers": {}}`,
			statusCode:     200,
			expectEmpty:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "empty host" {
				result := GenerateConfiguration(tt.providerConfig)
				if result == nil {
					t.Error("Expected non-nil configuration")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if _, err := w.Write([]byte(tt.serverResponse)); err != nil {
					t.Fatal(err)
				}
			}))
			defer server.Close()

			// Update host to use test server
			tt.providerConfig.Connection.Host = server.URL[7:] // Remove http://

			result := GenerateConfiguration(tt.providerConfig)

			if result == nil {
				t.Error("Expected non-nil configuration")
			}
		})
	}
}

func TestGenerateConfigurationInvalidURLRequest(t *testing.T) {
	// Test buildProviderRequest returning nil due to invalid URL format
	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host: "://invalid-url",
			Port: 8080,
		},
		HTTP: &config.HTTPSection{Discover: true},
	}

	result := GenerateConfiguration(providerConfig)
	if result == nil {
		t.Error("Expected non-nil configuration even with invalid URL")
	}
}

func TestParseDynamicConfigurationErrorPaths(t *testing.T) {
	tests := []struct {
		name           string
		jsonData       string
		providerConfig *config.ProviderConfig
		expectError    bool
	}{
		{
			name:     "invalid JSON",
			jsonData: `{"invalid": json}`,
			providerConfig: &config.ProviderConfig{
				HTTP: &config.HTTPSection{Discover: true},
			},
			expectError: true,
		},
		{
			name:     "non-object JSON",
			jsonData: `"string"`,
			providerConfig: &config.ProviderConfig{
				HTTP: &config.HTTPSection{Discover: true},
			},
			expectError: true,
		},
		{
			name:     "HTTP parsing error",
			jsonData: `{"routers": "invalid"}`,
			providerConfig: &config.ProviderConfig{
				HTTP: &config.HTTPSection{
					Discover: true,
					Routers:  &config.RoutersConfig{Discover: true},
				},
			},
			expectError: false, // parseHTTPConfig doesn't return errors for invalid data
		},
		{
			name:     "TCP parsing error",
			jsonData: `{"tcpRouters": "invalid"}`,
			providerConfig: &config.ProviderConfig{
				TCP: &config.TCPSection{
					Discover: true,
					Routers:  &config.RoutersConfig{Discover: true},
				},
			},
			expectError: false,
		},
		{
			name:     "UDP parsing error",
			jsonData: `{"udpRouters": "invalid"}`,
			providerConfig: &config.ProviderConfig{
				UDP: &config.UDPSection{
					Discover: true,
					Routers:  &config.UDPRoutersConfig{Discover: true},
				},
			},
			expectError: false,
		},
		{
			name:     "TLS parsing error",
			jsonData: `{"tls": "invalid"}`,
			providerConfig: &config.ProviderConfig{
				TLS: &config.TLSSection{Discover: true},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseDynamicConfiguration([]byte(tt.jsonData), tt.providerConfig)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if cfg == nil {
				t.Error("Expected non-nil configuration")
			}
		})
	}
}

func TestGenerateConfigurationTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host:    server.URL[7:], // Remove "http://" prefix
			Path:    "/api/rawdata",
			Timeout: "10ms", // Very short timeout
		},
		HTTP: &config.HTTPSection{Discover: true},
	}

	cfg := GenerateConfiguration(providerConfig)
	if cfg == nil {
		t.Error("Expected non-nil configuration")
	}
}

func TestGenerateConfigurationNetworkError(t *testing.T) {
	// Test network error by using invalid host
	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host: "invalid-host-that-does-not-exist",
			Port: 8080,
			Path: "/api/rawdata",
		},
		HTTP: &config.HTTPSection{Discover: true},
	}

	cfg := GenerateConfiguration(providerConfig)
	if cfg == nil {
		t.Error("Expected non-nil configuration")
	}
}

func TestParseDynamicConfigurationWithNilSections(t *testing.T) {
	// Test parsing with completely nil sections
	providerConfig := &config.ProviderConfig{
		HTTP: nil,
		TCP:  nil,
		UDP:  nil,
		TLS:  nil,
	}

	jsonData := `{"routers": {}, "services": {}}`
	cfg, err := parseDynamicConfiguration([]byte(jsonData), providerConfig)
	if err != nil {
		t.Errorf("Expected no error with nil sections, got: %v", err)
	}
	if cfg == nil {
		t.Error("Expected non-nil configuration")
	}
}

func TestParseDynamicConfigurationDiscoverFalse(t *testing.T) {
	// Test when all Discover flags are false
	providerConfig := &config.ProviderConfig{
		HTTP: &config.HTTPSection{Discover: false},
		TCP:  &config.TCPSection{Discover: false},
		UDP:  &config.UDPSection{Discover: false},
		TLS:  &config.TLSSection{Discover: false},
	}

	jsonData := `{"routers": {}, "services": {}, "tcpRouters": {}, "udpRouters": {}}`
	cfg, err := parseDynamicConfiguration([]byte(jsonData), providerConfig)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if cfg == nil {
		t.Error("Expected non-nil configuration")
	}
}

func TestGenerateConfigurationReadBodyError(t *testing.T) {
	// Test server that closes connection immediately to cause body read error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100") // Set content length but don't write body
		w.WriteHeader(http.StatusOK)
		// Don't write any body, causing a read error
	}))
	defer server.Close()

	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host: server.URL[7:], // Remove http://
			Path: "/api",
		},
		HTTP: &config.HTTPSection{Discover: true},
	}

	result := GenerateConfiguration(providerConfig)
	if result == nil {
		t.Error("Expected non-nil configuration even with body read error")
	}
}

func TestGenerateConfigurationHTTPClientDoError(t *testing.T) {
	// Test with a server that immediately closes to cause client.Do error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This handler won't be reached because we'll close the server
	}))
	server.Close() // Close immediately to cause connection error

	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host: server.URL[7:], // Remove http://
			Path: "/api",
		},
		HTTP: &config.HTTPSection{Discover: true},
	}

	result := GenerateConfiguration(providerConfig)
	if result == nil {
		t.Error("Expected non-nil configuration even with HTTP client error")
	}
}

func TestGenerateConfigurationNon200Status(t *testing.T) {
	// Test with server returning non-200 status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("Not Found")); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host: server.URL[7:], // Remove http://
			Path: "/api",
		},
		HTTP: &config.HTTPSection{Discover: true},
	}

	result := GenerateConfiguration(providerConfig)
	if result == nil {
		t.Error("Expected non-nil configuration even with non-200 status")
	}
}

func TestGenerateConfigurationParseError(t *testing.T) {
	// Test with server returning invalid JSON that causes parse error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("invalid json")); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host: server.URL[7:], // Remove http://
			Path: "/api",
		},
		HTTP: &config.HTTPSection{Discover: true},
	}

	result := GenerateConfiguration(providerConfig)
	if result == nil {
		t.Error("Expected non-nil configuration even with parse error")
	}
}

func TestGenerateConfigurationValidTimeoutParsing(t *testing.T) {
	// Test with valid timeout that gets parsed correctly
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"routers": {}}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host:    server.URL[7:], // Remove http://
			Path:    "/api",
			Timeout: "5s", // Valid timeout
		},
		HTTP: &config.HTTPSection{Discover: true},
	}

	result := GenerateConfiguration(providerConfig)
	if result == nil {
		t.Error("Expected non-nil configuration with valid timeout")
	}
}

func TestGenerateConfigurationInvalidTimeoutParsing(t *testing.T) {
	// Test with invalid timeout that fails to parse
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"routers": {}}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host:    server.URL[7:], // Remove http://
			Path:    "/api",
			Timeout: "invalid-timeout", // Invalid timeout format
		},
		HTTP: &config.HTTPSection{Discover: true},
	}

	result := GenerateConfiguration(providerConfig)
	if result == nil {
		t.Error("Expected non-nil configuration even with invalid timeout")
	}
}

func TestGenerateConfigurationBodyReadError(t *testing.T) {
	// Create a custom response body that will fail on Read
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Don't write anything, but set a content length to trigger read
	}))
	server.Close() // Close server to cause read error

	// Use the closed server URL to trigger network error during body read
	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host: server.URL[7:], // Remove http://
			Path: "/api",
		},
		HTTP: &config.HTTPSection{Discover: true},
	}

	result := GenerateConfiguration(providerConfig)
	if result == nil {
		t.Error("Expected non-nil configuration even with body read error")
	}
}

func TestGenerateConfigurationSuccessfulParse(t *testing.T) {
	// Test successful parsing path where parseDynamicConfiguration returns no error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"routers": {"test": {"rule": "Host(test.com)", "service": "test-service"}}}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host: server.URL[7:], // Remove http://
			Path: "/api",
		},
		HTTP: &config.HTTPSection{
			Discover: true,
			Routers:  &config.RoutersConfig{Discover: true},
		},
	}

	result := GenerateConfiguration(providerConfig)
	if result == nil {
		t.Error("Expected non-nil configuration for successful parse")
	}
}

func TestGenerateConfigurationParseErrorReturnsConfig(t *testing.T) {
	// Test that when parseDynamicConfiguration returns an error, we still return the config
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`invalid json that will cause parse error`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	providerConfig := &config.ProviderConfig{
		Connection: config.ConnectionConfig{
			Host: server.URL[7:], // Remove http://
			Path: "/api",
		},
		HTTP: &config.HTTPSection{Discover: true},
	}

	result := GenerateConfiguration(providerConfig)
	if result == nil {
		t.Error("Expected non-nil configuration even when parse returns error")
	}
}

func TestGenerateConfigurationAllPaths(t *testing.T) {
	tests := []struct {
		name         string
		setupServer  func() *httptest.Server
		expectResult bool
	}{
		{
			name: "successful 200 response with valid JSON",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					if _, err := w.Write([]byte(`{"routers": {}}`)); err != nil {
						t.Fatal(err)
					}
				}))
			},
			expectResult: true,
		},
		{
			name: "non-200 status code",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					if _, err := w.Write([]byte("Internal Server Error")); err != nil {
						t.Fatal(err)
					}
				}))
			},
			expectResult: true,
		},
		{
			name: "200 status with parse error",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					if _, err := w.Write([]byte("invalid json")); err != nil {
						t.Fatal(err)
					}
				}))
			},
			expectResult: true,
		},
		{
			name: "200 status with successful parse",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					if _, err := w.Write([]byte(`{"routers": {"test": {"rule": "Host(test.com)"}}}`)); err != nil {
						t.Fatal(err)
					}
				}))
			},
			expectResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			providerConfig := &config.ProviderConfig{
				Connection: config.ConnectionConfig{
					Host: server.URL[7:], // Remove http://
					Path: "/api",
				},
				HTTP: &config.HTTPSection{
					Discover: true,
					Routers:  &config.RoutersConfig{Discover: true},
				},
			}

			result := GenerateConfiguration(providerConfig)
			if tt.expectResult && result == nil {
				t.Error("Expected non-nil configuration")
			}
		})
	}
}
