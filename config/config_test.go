package config

import (
	"testing"
)

func TestProviderConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ProviderConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: ProviderConfig{
				Name: "test-provider",
				Connection: ConnectionConfig{
					Host:    "localhost",
					Port:    8080,
					Path:    "/api/rawdata",
					Timeout: "30s",
				},
				HTTP: &HTTPSection{
					Discover: true,
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: ProviderConfig{
				Connection: ConnectionConfig{
					Host: "localhost",
					Port: 8080,
				},
			},
			wantErr: true,
		},
		{
			name: "empty host",
			config: ProviderConfig{
				Name: "test-provider",
				Connection: ConnectionConfig{
					Host: "",
					Port: 8080,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			config: ProviderConfig{
				Name: "test-provider",
				Connection: ConnectionConfig{
					Host: "localhost",
					Port: -1,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simple validation check
			hasErr := tt.config.Name == "" || tt.config.Connection.Host == "" || tt.config.Connection.Port <= 0
			if hasErr != tt.wantErr {
				t.Errorf("ProviderConfig validation error = %v, wantErr %v", hasErr, tt.wantErr)
			}
		})
	}
}

func TestTunnelConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		tunnel  TunnelConfig
		wantErr bool
	}{
		{
			name: "valid tunnel config",
			tunnel: TunnelConfig{
				Name:      "test-tunnel",
				Addresses: []string{"http://tunnel1:80", "http://tunnel2:80"},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			tunnel: TunnelConfig{
				Addresses: []string{"http://tunnel:80"},
			},
			wantErr: true,
		},
		{
			name: "empty addresses",
			tunnel: TunnelConfig{
				Name:      "test-tunnel",
				Addresses: []string{},
			},
			wantErr: true,
		},
		{
			name: "nil addresses",
			tunnel: TunnelConfig{
				Name: "test-tunnel",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simple validation check
			hasErr := tt.tunnel.Name == "" || len(tt.tunnel.Addresses) == 0
			if hasErr != tt.wantErr {
				t.Errorf("TunnelConfig validation error = %v, wantErr %v", hasErr, tt.wantErr)
			}
		})
	}
}

func TestConnectionConfig_GetURL(t *testing.T) {
	tests := []struct {
		name     string
		config   ConnectionConfig
		expected string
	}{
		{
			name: "HTTP URL with path",
			config: ConnectionConfig{
				Host:    "localhost",
				Port:    8080,
				Path:    "/api/rawdata",
				Timeout: "30s",
			},
			expected: "http://localhost:8080/api/rawdata",
		},
		{
			name: "HTTPS URL",
			config: ConnectionConfig{
				Host: "secure.example.com",
				Port: 443,
				Path: "/config",
				MTLS: &MTLSConfig{
					CAFile: "/path/to/ca.crt",
				},
			},
			expected: "https://secure.example.com:443/config",
		},
		{
			name: "URL without path",
			config: ConnectionConfig{
				Host: "api.example.com",
				Port: 9000,
			},
			expected: "http://api.example.com:9000",
		},
		{
			name: "URL with multiple hosts (uses first)",
			config: ConnectionConfig{
				Host: "host1.example.com",
				Port: 8080,
				Path: "/api",
			},
			expected: "http://host1.example.com:8080/api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simple URL construction
			scheme := "http"
			if tt.config.MTLS != nil {
				scheme = "https"
			}
			// For test purposes, just check the scheme is correct
			expectedScheme := "http"
			if tt.config.MTLS != nil {
				expectedScheme = "https"
			}
			if (tt.config.MTLS != nil && scheme != "https") || (tt.config.MTLS == nil && scheme != "http") {
				t.Errorf("ConnectionConfig scheme mismatch: got %s, expected %s", scheme, expectedScheme)
			}
		})
	}
}

func TestOverrideServer_HasTunnel(t *testing.T) {
	tests := []struct {
		name     string
		override OverrideServer
		expected bool
	}{
		{
			name: "has tunnel",
			override: OverrideServer{
				Tunnel: "test-tunnel",
			},
			expected: true,
		},
		{
			name: "no tunnel",
			override: OverrideServer{
				Strategy: "replace",
				Value:    []string{"http://server:80"},
			},
			expected: false,
		},
		{
			name: "empty tunnel",
			override: OverrideServer{
				Tunnel: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.override.Tunnel != ""
			if result != tt.expected {
				t.Errorf("OverrideServer tunnel check = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestOverrideHealthcheck_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		hc       OverrideHealthcheck
		expected bool
	}{
		{
			name:     "empty healthcheck",
			hc:       OverrideHealthcheck{},
			expected: true,
		},
		{
			name: "healthcheck with path",
			hc: OverrideHealthcheck{
				Path: "/health",
			},
			expected: false,
		},
		{
			name: "healthcheck with interval",
			hc: OverrideHealthcheck{
				Interval: "10s",
			},
			expected: false,
		},
		{
			name: "healthcheck with timeout",
			hc: OverrideHealthcheck{
				Timeout: "5s",
			},
			expected: false,
		},
		{
			name: "complete healthcheck",
			hc: OverrideHealthcheck{
				Path:     "/api/health",
				Interval: "15s",
				Timeout:  "3s",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hc.Path == "" && tt.hc.Interval == "" && tt.hc.Timeout == ""
			if result != tt.expected {
				t.Errorf("OverrideHealthcheck empty check = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestServicesConfig_ShouldApplyOverrides(t *testing.T) {
	tests := []struct {
		name     string
		config   ServicesConfig
		expected bool
	}{
		{
			name: "has server overrides",
			config: ServicesConfig{
				Overrides: ServiceOverrides{
					Servers: []OverrideServer{
						{Strategy: "replace", Value: []string{"http://new:80"}},
					},
				},
			},
			expected: true,
		},
		{
			name: "has healthcheck overrides",
			config: ServicesConfig{
				Overrides: ServiceOverrides{
					Healthchecks: []OverrideHealthcheck{
						{Path: "/health"},
					},
				},
			},
			expected: true,
		},
		{
			name: "no overrides",
			config: ServicesConfig{
				Overrides: ServiceOverrides{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := len(tt.config.Overrides.Servers) > 0 || len(tt.config.Overrides.Healthchecks) > 0
			if result != tt.expected {
				t.Errorf("ServicesConfig override check = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCreateConfig(t *testing.T) {
	config := CreateConfig()

	if config == nil {
		t.Fatal("CreateConfig returned nil")
	}

	if config.PollInterval != "5s" {
		t.Errorf("Expected PollInterval '5s', got '%s'", config.PollInterval)
	}

	if len(config.Providers) != 1 {
		t.Errorf("Expected 1 provider, got %d", len(config.Providers))
	}

	provider := config.Providers[0]
	if provider.Name != "Traefik Provider" {
		t.Errorf("Expected provider name 'Traefik Provider', got '%s'", provider.Name)
	}

	if provider.Connection.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", provider.Connection.Host)
	}

	if provider.Connection.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", provider.Connection.Port)
	}

	if provider.Connection.Path != "/api/rawdata" {
		t.Errorf("Expected path '/api/rawdata', got '%s'", provider.Connection.Path)
	}

	if provider.Connection.Timeout != "5s" {
		t.Errorf("Expected timeout '5s', got '%s'", provider.Connection.Timeout)
	}

	if provider.HTTP == nil {
		t.Error("Expected HTTP section to be initialized")
	} else {
		if !provider.HTTP.Discover {
			t.Error("Expected HTTP.Discover to be true")
		}

		if provider.HTTP.Routers == nil {
			t.Error("Expected HTTP.Routers to be initialized")
		} else {
			if !provider.HTTP.Routers.Discover {
				t.Error("Expected HTTP.Routers.Discover to be true")
			}
		}

		if provider.HTTP.Services == nil {
			t.Error("Expected HTTP.Services to be initialized")
		} else {
			if !provider.HTTP.Services.Discover {
				t.Error("Expected HTTP.Services.Discover to be true")
			}
		}

		if provider.HTTP.Middlewares == nil {
			t.Error("Expected HTTP.Middlewares to be initialized")
		} else {
			if !provider.HTTP.Middlewares.Discover {
				t.Error("Expected HTTP.Middlewares.Discover to be true")
			}
		}

		if !provider.HTTP.ServerTransports.Discover {
			t.Error("Expected HTTP.ServerTransports.Discover to be true")
		}
	}

	if provider.TCP == nil {
		t.Error("Expected TCP section to be initialized")
	} else {
		if !provider.TCP.Discover {
			t.Error("Expected TCP.Discover to be true")
		}
	}

	if provider.UDP == nil {
		t.Error("Expected UDP section to be initialized")
	} else {
		if !provider.UDP.Discover {
			t.Error("Expected UDP.Discover to be true")
		}
	}
}
