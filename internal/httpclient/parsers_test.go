package httpclient

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func TestParseHTTPConfigMarshalErrors(t *testing.T) {
	httpConfig := &dynamic.HTTPConfiguration{
		Routers:     make(map[string]*dynamic.Router),
		Services:    make(map[string]*dynamic.Service),
		Middlewares: make(map[string]*dynamic.Middleware),
	}

	providerConfig := &config.HTTPSection{
		Routers: &config.RoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				make(chan int), // Unmarshalable type
				map[string]interface{}{
					"name": "test-router",
					"rule": "Host(`test.com`)",
				},
			},
		},
		Services: &config.ServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				make(chan int), // Unmarshalable type
			},
		},
		Middlewares: &config.MiddlewaresConfig{
			Discover: true,
			ExtraMiddlewares: []interface{}{
				make(chan int), // Unmarshalable type
			},
		},
	}

	raw := map[string]interface{}{}
	err := parseHTTPConfig(raw, httpConfig, providerConfig, nil)
	if err != nil {
		t.Errorf("parseHTTPConfig should not return error: %v", err)
	}

	if len(httpConfig.Routers) != 1 {
		t.Errorf("Expected 1 router, got %d", len(httpConfig.Routers))
	}
}

func TestParseHTTPConfig(t *testing.T) {
	tests := []struct {
		name           string
		raw            map[string]interface{}
		httpConfig     *dynamic.HTTPConfiguration
		providerConfig *config.HTTPSection
		tunnels        []config.TunnelConfig
		expectError    bool
	}{
		{
			name: "parse routers when discover enabled",
			raw: map[string]interface{}{
				"routers": map[string]interface{}{
					"test-router": map[string]interface{}{
						"rule":    "Host(`example.com`)",
						"service": "test-service",
					},
				},
			},
			httpConfig: &dynamic.HTTPConfiguration{
				Routers: make(map[string]*dynamic.Router),
			},
			providerConfig: &config.HTTPSection{
				Routers: &config.RoutersConfig{
					Discover: true,
				},
				Services: &config.ServicesConfig{
					Discover: false,
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: false,
				},
			},
		},
		{
			name: "parse services when discover enabled",
			raw: map[string]interface{}{
				"services": map[string]interface{}{
					"test-service": map[string]interface{}{
						"loadBalancer": map[string]interface{}{
							"servers": []interface{}{
								map[string]interface{}{"url": "http://backend:8080"},
							},
						},
					},
				},
			},
			httpConfig: &dynamic.HTTPConfiguration{
				Services: make(map[string]*dynamic.Service),
			},
			providerConfig: &config.HTTPSection{
				Routers: &config.RoutersConfig{
					Discover: false,
				},
				Services: &config.ServicesConfig{
					Discover: true,
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: false,
				},
			},
		},
		{
			name: "parse middlewares when discover enabled",
			raw: map[string]interface{}{
				"middlewares": map[string]interface{}{
					"test-middleware": map[string]interface{}{
						"stripPrefix": map[string]interface{}{
							"prefixes": []string{"/api"},
						},
					},
				},
			},
			httpConfig: &dynamic.HTTPConfiguration{
				Middlewares: make(map[string]*dynamic.Middleware),
			},
			providerConfig: &config.HTTPSection{
				Routers: &config.RoutersConfig{
					Discover: false,
				},
				Services: &config.ServicesConfig{
					Discover: false,
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: true,
				},
			},
		},
		{
			name: "add extra routes",
			raw:  map[string]interface{}{},
			httpConfig: &dynamic.HTTPConfiguration{
				Routers: make(map[string]*dynamic.Router),
			},
			providerConfig: &config.HTTPSection{
				Routers: &config.RoutersConfig{
					Discover: true,
					ExtraRoutes: []interface{}{
						map[string]interface{}{
							"name":    "extra-router",
							"rule":    "Host(`extra.com`)",
							"service": "extra-service",
						},
					},
				},
				Services: &config.ServicesConfig{
					Discover: false,
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: false,
				},
			},
		},
		{
			name: "add extra services",
			raw:  map[string]interface{}{},
			httpConfig: &dynamic.HTTPConfiguration{
				Services: make(map[string]*dynamic.Service),
			},
			providerConfig: &config.HTTPSection{
				Routers: &config.RoutersConfig{
					Discover: false,
				},
				Services: &config.ServicesConfig{
					Discover: true,
					ExtraServices: []interface{}{
						map[string]interface{}{
							"name": "extra-service",
							"loadBalancer": map[string]interface{}{
								"servers": []interface{}{
									map[string]interface{}{"url": "http://extra:8080"},
								},
							},
						},
					},
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: false,
				},
			},
		},
		{
			name: "add extra middlewares",
			raw:  map[string]interface{}{},
			httpConfig: &dynamic.HTTPConfiguration{
				Middlewares: make(map[string]*dynamic.Middleware),
			},
			providerConfig: &config.HTTPSection{
				Routers: &config.RoutersConfig{
					Discover: false,
				},
				Services: &config.ServicesConfig{
					Discover: false,
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: true,
					ExtraMiddlewares: []interface{}{
						map[string]interface{}{
							"name": "extra-middleware",
							"stripPrefix": map[string]interface{}{
								"prefixes": []string{"/extra"},
							},
						},
					},
				},
			},
		},
		{
			name: "invalid extra route without name",
			raw:  map[string]interface{}{},
			httpConfig: &dynamic.HTTPConfiguration{
				Routers: make(map[string]*dynamic.Router),
			},
			providerConfig: &config.HTTPSection{
				Routers: &config.RoutersConfig{
					Discover: true,
					ExtraRoutes: []interface{}{
						map[string]interface{}{
							"rule":    "Host(`invalid.com`)",
							"service": "invalid-service",
							// missing name
						},
					},
				},
				Services: &config.ServicesConfig{
					Discover: false,
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: false,
				},
			},
		},
		{
			name: "invalid extra route with unmarshalable data",
			raw:  map[string]interface{}{},
			httpConfig: &dynamic.HTTPConfiguration{
				Routers: make(map[string]*dynamic.Router),
			},
			providerConfig: &config.HTTPSection{
				Routers: &config.RoutersConfig{
					Discover: true,
					ExtraRoutes: []interface{}{
						func() {}, // This will cause json.Marshal to fail
					},
				},
				Services: &config.ServicesConfig{
					Discover: false,
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: false,
				},
			},
		},
		{
			name:       "no discover sections enabled",
			raw:        map[string]interface{}{},
			httpConfig: &dynamic.HTTPConfiguration{},
			providerConfig: &config.HTTPSection{
				Routers:     &config.RoutersConfig{Discover: false},
				Services:    &config.ServicesConfig{Discover: false},
				Middlewares: &config.MiddlewaresConfig{Discover: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseHTTPConfig(tt.raw, tt.httpConfig, tt.providerConfig, tt.tunnels)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestParseTCPConfig(t *testing.T) {
	tests := []struct {
		name           string
		raw            map[string]interface{}
		tcpConfig      *dynamic.TCPConfiguration
		providerConfig *config.TCPSection
		tunnels        []config.TunnelConfig
		expectError    bool
	}{
		{
			name: "parse tcp routers when discover enabled",
			raw: map[string]interface{}{
				"tcpRouters": map[string]interface{}{
					"tcp-router": map[string]interface{}{
						"rule":    "HostSNI(`tcp.example.com`)",
						"service": "tcp-service",
					},
				},
			},
			tcpConfig: &dynamic.TCPConfiguration{
				Routers: make(map[string]*dynamic.TCPRouter),
			},
			providerConfig: &config.TCPSection{
				Routers: &config.RoutersConfig{
					Discover: true,
				},
				Services: &config.ServicesConfig{
					Discover: false,
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: false,
				},
			},
		},
		{
			name: "parse tcp services when discover enabled",
			raw: map[string]interface{}{
				"tcpServices": map[string]interface{}{
					"tcp-service": map[string]interface{}{
						"loadBalancer": map[string]interface{}{
							"servers": []interface{}{
								map[string]interface{}{"address": "backend:8081"},
							},
						},
					},
				},
			},
			tcpConfig: &dynamic.TCPConfiguration{
				Services: make(map[string]*dynamic.TCPService),
			},
			providerConfig: &config.TCPSection{
				Routers: &config.RoutersConfig{
					Discover: false,
				},
				Services: &config.ServicesConfig{
					Discover: true,
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: false,
				},
			},
		},
		{
			name: "parse tcp middlewares when discover enabled",
			raw: map[string]interface{}{
				"tcpMiddlewares": map[string]interface{}{
					"tcp-middleware": map[string]interface{}{
						"ipWhiteList": map[string]interface{}{
							"sourceRange": []string{"192.168.1.0/24"},
						},
					},
				},
			},
			tcpConfig: &dynamic.TCPConfiguration{
				Middlewares: make(map[string]*dynamic.TCPMiddleware),
			},
			providerConfig: &config.TCPSection{
				Routers: &config.RoutersConfig{
					Discover: false,
				},
				Services: &config.ServicesConfig{
					Discover: false,
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: true,
				},
			},
		},
		{
			name: "add extra tcp routes",
			raw:  map[string]interface{}{},
			tcpConfig: &dynamic.TCPConfiguration{
				Routers: make(map[string]*dynamic.TCPRouter),
			},
			providerConfig: &config.TCPSection{
				Routers: &config.RoutersConfig{
					Discover: true,
					ExtraRoutes: []interface{}{
						map[string]interface{}{
							"name":    "extra-tcp-router",
							"rule":    "HostSNI(`extra-tcp.com`)",
							"service": "extra-tcp-service",
						},
					},
				},
				Services: &config.ServicesConfig{
					Discover: false,
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: false,
				},
			},
		},
		{
			name: "add extra tcp services",
			raw:  map[string]interface{}{},
			tcpConfig: &dynamic.TCPConfiguration{
				Services: make(map[string]*dynamic.TCPService),
			},
			providerConfig: &config.TCPSection{
				Routers: &config.RoutersConfig{
					Discover: false,
				},
				Services: &config.ServicesConfig{
					Discover: true,
					ExtraServices: []interface{}{
						map[string]interface{}{
							"name": "extra-tcp-service",
							"loadBalancer": map[string]interface{}{
								"servers": []interface{}{
									map[string]interface{}{"address": "extra-tcp:8081"},
								},
							},
						},
					},
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: false,
				},
			},
		},
		{
			name: "add extra tcp middlewares",
			raw:  map[string]interface{}{},
			tcpConfig: &dynamic.TCPConfiguration{
				Middlewares: make(map[string]*dynamic.TCPMiddleware),
			},
			providerConfig: &config.TCPSection{
				Routers: &config.RoutersConfig{
					Discover: false,
				},
				Services: &config.ServicesConfig{
					Discover: false,
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: true,
					ExtraMiddlewares: []interface{}{
						map[string]interface{}{
							"name": "extra-tcp-middleware",
							"ipWhiteList": map[string]interface{}{
								"sourceRange": []string{"10.0.0.0/8"},
							},
						},
					},
				},
			},
		},
		{
			name: "invalid extra tcp route without name",
			raw:  map[string]interface{}{},
			tcpConfig: &dynamic.TCPConfiguration{
				Routers: make(map[string]*dynamic.TCPRouter),
			},
			providerConfig: &config.TCPSection{
				Routers: &config.RoutersConfig{
					Discover: true,
					ExtraRoutes: []interface{}{
						map[string]interface{}{
							"rule":    "HostSNI(`invalid-tcp.com`)",
							"service": "invalid-tcp-service",
							// missing name
						},
					},
				},
				Services: &config.ServicesConfig{
					Discover: false,
				},
				Middlewares: &config.MiddlewaresConfig{
					Discover: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseTCPConfig(tt.raw, tt.tcpConfig, tt.providerConfig, tt.tunnels)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestParseUDPConfig(t *testing.T) {
	tests := []struct {
		name           string
		raw            map[string]interface{}
		udpConfig      *dynamic.UDPConfiguration
		providerConfig *config.UDPSection
		tunnels        []config.TunnelConfig
		expectError    bool
	}{
		{
			name: "parse udp routers when discover enabled",
			raw: map[string]interface{}{
				"udpRouters": map[string]interface{}{
					"udp-router": map[string]interface{}{
						"service": "udp-service",
					},
				},
			},
			udpConfig: &dynamic.UDPConfiguration{
				Routers: make(map[string]*dynamic.UDPRouter),
			},
			providerConfig: &config.UDPSection{
				Routers: &config.UDPRoutersConfig{
					Discover: true,
				},
				Services: &config.UDPServicesConfig{
					Discover: false,
				},
			},
		},
		{
			name: "parse udp services when discover enabled",
			raw: map[string]interface{}{
				"udpServices": map[string]interface{}{
					"udp-service": map[string]interface{}{
						"loadBalancer": map[string]interface{}{
							"servers": []interface{}{
								map[string]interface{}{"address": "backend:8082"},
							},
						},
					},
				},
			},
			udpConfig: &dynamic.UDPConfiguration{
				Services: make(map[string]*dynamic.UDPService),
			},
			providerConfig: &config.UDPSection{
				Routers: &config.UDPRoutersConfig{
					Discover: false,
				},
				Services: &config.UDPServicesConfig{
					Discover: true,
				},
			},
		},
		{
			name: "add extra udp routes",
			raw:  map[string]interface{}{},
			udpConfig: &dynamic.UDPConfiguration{
				Routers: make(map[string]*dynamic.UDPRouter),
			},
			providerConfig: &config.UDPSection{
				Routers: &config.UDPRoutersConfig{
					Discover: true,
					ExtraRoutes: []interface{}{
						map[string]interface{}{
							"name":    "extra-udp-router",
							"service": "extra-udp-service",
						},
					},
				},
				Services: &config.UDPServicesConfig{
					Discover: false,
				},
			},
		},
		{
			name: "add extra udp services",
			raw:  map[string]interface{}{},
			udpConfig: &dynamic.UDPConfiguration{
				Services: make(map[string]*dynamic.UDPService),
			},
			providerConfig: &config.UDPSection{
				Routers: &config.UDPRoutersConfig{
					Discover: false,
				},
				Services: &config.UDPServicesConfig{
					Discover: true,
					ExtraServices: []interface{}{
						map[string]interface{}{
							"name": "extra-udp-service",
							"loadBalancer": map[string]interface{}{
								"servers": []interface{}{
									map[string]interface{}{"address": "extra-udp:8082"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "invalid extra udp route without name",
			raw:  map[string]interface{}{},
			udpConfig: &dynamic.UDPConfiguration{
				Routers: make(map[string]*dynamic.UDPRouter),
			},
			providerConfig: &config.UDPSection{
				Routers: &config.UDPRoutersConfig{
					Discover: true,
					ExtraRoutes: []interface{}{
						map[string]interface{}{
							"service": "invalid-udp-service",
							// missing name
						},
					},
				},
				Services: &config.UDPServicesConfig{
					Discover: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseUDPConfig(tt.raw, tt.udpConfig, tt.providerConfig, tt.tunnels)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestParseTLSConfig(t *testing.T) {
	tests := []struct {
		name           string
		raw            map[string]interface{}
		tlsConfig      *dynamic.TLSConfiguration
		providerConfig *config.TLSSection
		expectError    bool
	}{
		{
			name: "parse tls certificates",
			raw: map[string]interface{}{
				"tlsCertificates": []interface{}{
					map[string]interface{}{
						"certFile": "/path/to/cert.pem",
						"keyFile":  "/path/to/key.pem",
					},
				},
			},
			tlsConfig:      &dynamic.TLSConfiguration{},
			providerConfig: &config.TLSSection{Discover: true},
		},
		{
			name: "parse tls options",
			raw: map[string]interface{}{
				"tlsOptions": map[string]interface{}{
					"default": map[string]interface{}{
						"minVersion": "VersionTLS12",
					},
				},
			},
			tlsConfig:      &dynamic.TLSConfiguration{},
			providerConfig: &config.TLSSection{Discover: true},
		},
		{
			name: "parse tls stores",
			raw: map[string]interface{}{
				"tlsStores": map[string]interface{}{
					"default": map[string]interface{}{
						"defaultCertificate": map[string]interface{}{
							"certFile": "/path/to/default-cert.pem",
							"keyFile":  "/path/to/default-key.pem",
						},
					},
				},
			},
			tlsConfig:      &dynamic.TLSConfiguration{},
			providerConfig: &config.TLSSection{Discover: true},
		},
		{
			name: "parse all tls sections",
			raw: map[string]interface{}{
				"tlsCertificates": []interface{}{
					map[string]interface{}{
						"certFile": "/path/to/cert.pem",
						"keyFile":  "/path/to/key.pem",
					},
				},
				"tlsOptions": map[string]interface{}{
					"default": map[string]interface{}{
						"minVersion": "VersionTLS12",
					},
				},
				"tlsStores": map[string]interface{}{
					"default": map[string]interface{}{
						"defaultCertificate": map[string]interface{}{
							"certFile": "/path/to/default-cert.pem",
							"keyFile":  "/path/to/default-key.pem",
						},
					},
				},
			},
			tlsConfig:      &dynamic.TLSConfiguration{},
			providerConfig: &config.TLSSection{Discover: true},
		},
		{
			name:           "empty raw data",
			raw:            map[string]interface{}{},
			tlsConfig:      &dynamic.TLSConfiguration{},
			providerConfig: &config.TLSSection{Discover: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseTLSConfig(tt.raw, tt.tlsConfig, tt.providerConfig)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestParseHTTPConfigJSONMarshalErrors(t *testing.T) {
	httpConfig := &dynamic.HTTPConfiguration{
		Routers:     make(map[string]*dynamic.Router),
		Services:    make(map[string]*dynamic.Service),
		Middlewares: make(map[string]*dynamic.Middleware),
	}

	providerConfig := &config.HTTPSection{
		Routers: &config.RoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				make(chan int), // Unmarshalable type - triggers json.Marshal error
				map[string]interface{}{
					"name": "valid-router",
					"rule": "Host(`test.com`)",
				},
			},
		},
		Services: &config.ServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				make(chan int), // Unmarshalable type
				map[string]interface{}{
					"name": "valid-service",
					"loadBalancer": map[string]interface{}{
						"servers": []interface{}{
							map[string]interface{}{"url": "http://test:8080"},
						},
					},
				},
			},
		},
		Middlewares: &config.MiddlewaresConfig{
			Discover: true,
			ExtraMiddlewares: []interface{}{
				make(chan int), // Unmarshalable type
				map[string]interface{}{
					"name": "valid-middleware",
					"basicAuth": map[string]interface{}{
						"users": []string{"test:password"},
					},
				},
			},
		},
	}

	raw := map[string]interface{}{}
	err := parseHTTPConfig(raw, httpConfig, providerConfig, nil)
	if err != nil {
		t.Errorf("parseHTTPConfig should not return error: %v", err)
	}

	// Valid entries should be added despite marshal errors
	if len(httpConfig.Routers) != 1 {
		t.Errorf("Expected 1 router, got %d", len(httpConfig.Routers))
	}
	if len(httpConfig.Services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(httpConfig.Services))
	}
	if len(httpConfig.Middlewares) != 1 {
		t.Errorf("Expected 1 middleware, got %d", len(httpConfig.Middlewares))
	}
}

func TestParseHTTPConfigJSONUnmarshalErrors(t *testing.T) {
	httpConfig := &dynamic.HTTPConfiguration{
		Routers:     make(map[string]*dynamic.Router),
		Services:    make(map[string]*dynamic.Service),
		Middlewares: make(map[string]*dynamic.Middleware),
	}

	providerConfig := &config.HTTPSection{
		Routers: &config.RoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				map[string]interface{}{
					"name": "test-router",
					"rule": make(chan int), // Invalid data that will cause unmarshal error
				},
			},
		},
		Services: &config.ServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				map[string]interface{}{
					"name":         "test-service",
					"loadBalancer": make(chan int), // Invalid data
				},
			},
		},
		Middlewares: &config.MiddlewaresConfig{
			Discover: true,
			ExtraMiddlewares: []interface{}{
				map[string]interface{}{
					"name":      "test-middleware",
					"basicAuth": make(chan int), // Invalid data
				},
			},
		},
	}

	raw := map[string]interface{}{}
	err := parseHTTPConfig(raw, httpConfig, providerConfig, nil)
	if err != nil {
		t.Errorf("parseHTTPConfig should not return error: %v", err)
	}

	// Items with unmarshal errors should not be added
	if len(httpConfig.Routers) != 0 {
		t.Errorf("Expected 0 routers, got %d", len(httpConfig.Routers))
	}
	if len(httpConfig.Services) != 0 {
		t.Errorf("Expected 0 services, got %d", len(httpConfig.Services))
	}
	if len(httpConfig.Middlewares) != 0 {
		t.Errorf("Expected 0 middlewares, got %d", len(httpConfig.Middlewares))
	}
}

func TestParseTCPConfigJSONMarshalErrors(t *testing.T) {
	tcpConfig := &dynamic.TCPConfiguration{
		Routers:     make(map[string]*dynamic.TCPRouter),
		Services:    make(map[string]*dynamic.TCPService),
		Middlewares: make(map[string]*dynamic.TCPMiddleware),
	}

	providerConfig := &config.TCPSection{
		Discover: true,
		Routers: &config.RoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				make(chan int), // Unmarshalable type
				map[string]interface{}{
					"name": "valid-tcp-router",
					"rule": "HostSNI(`test.com`)",
				},
			},
		},
		Services: &config.ServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				make(chan int), // Unmarshalable type
			},
		},
		Middlewares: &config.MiddlewaresConfig{
			Discover: true,
			ExtraMiddlewares: []interface{}{
				make(chan int), // Unmarshalable type
			},
		},
	}

	raw := map[string]interface{}{}
	err := parseTCPConfig(raw, tcpConfig, providerConfig, nil)
	if err != nil {
		t.Errorf("parseTCPConfig should not return error: %v", err)
	}

	if len(tcpConfig.Routers) != 1 {
		t.Errorf("Expected 1 router, got %d", len(tcpConfig.Routers))
	}
}

func TestParseUDPConfigJSONMarshalErrors(t *testing.T) {
	udpConfig := &dynamic.UDPConfiguration{
		Routers:  make(map[string]*dynamic.UDPRouter),
		Services: make(map[string]*dynamic.UDPService),
	}

	providerConfig := &config.UDPSection{
		Discover: true,
		Routers: &config.UDPRoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				make(chan int), // Unmarshalable type
				map[string]interface{}{
					"name":    "valid-udp-router",
					"service": "test-service",
				},
			},
		},
		Services: &config.UDPServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				make(chan int), // Unmarshalable type
			},
		},
	}

	raw := map[string]interface{}{}
	err := parseUDPConfig(raw, udpConfig, providerConfig, nil)
	if err != nil {
		t.Errorf("parseUDPConfig should not return error: %v", err)
	}

	if len(udpConfig.Routers) != 1 {
		t.Errorf("Expected 1 router, got %d", len(udpConfig.Routers))
	}
}

func TestParseDynamicConfigurationNilSectionPointers(t *testing.T) {
	// Test when section pointers are nil but Discover is true
	providerConfig := &config.ProviderConfig{
		HTTP: &config.HTTPSection{
			Discover:    true,
			Routers:     nil, // This should not cause panic
			Services:    nil,
			Middlewares: nil,
		},
		TCP: &config.TCPSection{
			Discover:    true,
			Routers:     nil,
			Services:    nil,
			Middlewares: nil,
		},
		UDP: &config.UDPSection{
			Discover: true,
			Routers:  nil,
			Services: nil,
		},
		TLS: &config.TLSSection{
			Discover: true,
		},
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

func TestParseHTTPConfigExtraItemsWithoutName(t *testing.T) {
	httpConfig := &dynamic.HTTPConfiguration{
		Routers:     make(map[string]*dynamic.Router),
		Services:    make(map[string]*dynamic.Service),
		Middlewares: make(map[string]*dynamic.Middleware),
	}

	providerConfig := &config.HTTPSection{
		Routers: &config.RoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				map[string]interface{}{
					"rule":    "Host(`test.com`)",
					"service": "test-service",
					// Missing "name" field
				},
			},
		},
		Services: &config.ServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				map[string]interface{}{
					"loadBalancer": map[string]interface{}{
						"servers": []interface{}{
							map[string]interface{}{"url": "http://test:8080"},
						},
					},
					// Missing "name" field
				},
			},
		},
		Middlewares: &config.MiddlewaresConfig{
			Discover: true,
			ExtraMiddlewares: []interface{}{
				map[string]interface{}{
					"basicAuth": map[string]interface{}{
						"users": []string{"test:password"},
					},
					// Missing "name" field
				},
			},
		},
	}

	raw := map[string]interface{}{}
	err := parseHTTPConfig(raw, httpConfig, providerConfig, nil)
	if err != nil {
		t.Errorf("parseHTTPConfig should not return error: %v", err)
	}

	// Items without name should not be added
	if len(httpConfig.Routers) != 0 {
		t.Errorf("Expected 0 routers, got %d", len(httpConfig.Routers))
	}
	if len(httpConfig.Services) != 0 {
		t.Errorf("Expected 0 services, got %d", len(httpConfig.Services))
	}
	if len(httpConfig.Middlewares) != 0 {
		t.Errorf("Expected 0 middlewares, got %d", len(httpConfig.Middlewares))
	}
}

func TestParseTCPConfigExtraItemsWithoutName(t *testing.T) {
	tcpConfig := &dynamic.TCPConfiguration{
		Routers:     make(map[string]*dynamic.TCPRouter),
		Services:    make(map[string]*dynamic.TCPService),
		Middlewares: make(map[string]*dynamic.TCPMiddleware),
	}

	providerConfig := &config.TCPSection{
		Routers: &config.RoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				map[string]interface{}{
					"rule":    "HostSNI(`test.com`)",
					"service": "test-service",
					// Missing "name" field
				},
			},
		},
		Services: &config.ServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				map[string]interface{}{
					"loadBalancer": map[string]interface{}{
						"servers": []interface{}{
							map[string]interface{}{"address": "test:8081"},
						},
					},
					// Missing "name" field
				},
			},
		},
		Middlewares: &config.MiddlewaresConfig{
			Discover: true,
			ExtraMiddlewares: []interface{}{
				map[string]interface{}{
					"ipWhiteList": map[string]interface{}{
						"sourceRange": []string{"192.168.1.0/24"},
					},
					// Missing "name" field
				},
			},
		},
	}

	raw := map[string]interface{}{}
	err := parseTCPConfig(raw, tcpConfig, providerConfig, nil)
	if err != nil {
		t.Errorf("parseTCPConfig should not return error: %v", err)
	}

	// Items without name should not be added
	if len(tcpConfig.Routers) != 0 {
		t.Errorf("Expected 0 routers, got %d", len(tcpConfig.Routers))
	}
	if len(tcpConfig.Services) != 0 {
		t.Errorf("Expected 0 services, got %d", len(tcpConfig.Services))
	}
	if len(tcpConfig.Middlewares) != 0 {
		t.Errorf("Expected 0 middlewares, got %d", len(tcpConfig.Middlewares))
	}
}

func TestParseUDPConfigExtraItemsWithoutName(t *testing.T) {
	udpConfig := &dynamic.UDPConfiguration{
		Routers:  make(map[string]*dynamic.UDPRouter),
		Services: make(map[string]*dynamic.UDPService),
	}

	providerConfig := &config.UDPSection{
		Routers: &config.UDPRoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				map[string]interface{}{
					"service": "test-service",
					// Missing "name" field
				},
			},
		},
		Services: &config.UDPServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				map[string]interface{}{
					"loadBalancer": map[string]interface{}{
						"servers": []interface{}{
							map[string]interface{}{"address": "test:8082"},
						},
					},
					// Missing "name" field
				},
			},
		},
	}

	raw := map[string]interface{}{}
	err := parseUDPConfig(raw, udpConfig, providerConfig, nil)
	if err != nil {
		t.Errorf("parseUDPConfig should not return error: %v", err)
	}

	// Items without name should not be added
	if len(udpConfig.Routers) != 0 {
		t.Errorf("Expected 0 routers, got %d", len(udpConfig.Routers))
	}
	if len(udpConfig.Services) != 0 {
		t.Errorf("Expected 0 services, got %d", len(udpConfig.Services))
	}
}

func TestParseConfigMarshalErrors(t *testing.T) {
	// Test HTTP config marshal errors
	httpConfig := &dynamic.HTTPConfiguration{
		Routers:     make(map[string]*dynamic.Router),
		Services:    make(map[string]*dynamic.Service),
		Middlewares: make(map[string]*dynamic.Middleware),
	}

	httpProviderConfig := &config.HTTPSection{
		Routers: &config.RoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				// Channel cannot be marshaled - triggers marshal error continue
				make(chan int),
			},
		},
		Services: &config.ServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				// Channel cannot be marshaled - triggers marshal error continue
				make(chan int),
			},
		},
		Middlewares: &config.MiddlewaresConfig{
			Discover: true,
			ExtraMiddlewares: []interface{}{
				// Channel cannot be marshaled - triggers marshal error continue
				make(chan int),
			},
		},
	}

	err := parseHTTPConfig(map[string]interface{}{}, httpConfig, httpProviderConfig, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// All should be empty due to marshal errors
	if len(httpConfig.Routers) != 0 {
		t.Errorf("Expected 0 routers, got %d", len(httpConfig.Routers))
	}
	if len(httpConfig.Services) != 0 {
		t.Errorf("Expected 0 services, got %d", len(httpConfig.Services))
	}
	if len(httpConfig.Middlewares) != 0 {
		t.Errorf("Expected 0 middlewares, got %d", len(httpConfig.Middlewares))
	}

	// Test TCP config marshal errors
	tcpConfig := &dynamic.TCPConfiguration{
		Routers:     make(map[string]*dynamic.TCPRouter),
		Services:    make(map[string]*dynamic.TCPService),
		Middlewares: make(map[string]*dynamic.TCPMiddleware),
	}

	tcpProviderConfig := &config.TCPSection{
		Routers: &config.RoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				// Channel cannot be marshaled - triggers marshal error continue
				make(chan int),
			},
		},
		Services: &config.ServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				// Channel cannot be marshaled - triggers marshal error continue
				make(chan int),
			},
		},
		Middlewares: &config.MiddlewaresConfig{
			Discover: true,
			ExtraMiddlewares: []interface{}{
				// Channel cannot be marshaled - triggers marshal error continue
				make(chan int),
			},
		},
	}

	err = parseTCPConfig(map[string]interface{}{}, tcpConfig, tcpProviderConfig, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// All should be empty due to marshal errors
	if len(tcpConfig.Routers) != 0 {
		t.Errorf("Expected 0 routers, got %d", len(tcpConfig.Routers))
	}
	if len(tcpConfig.Services) != 0 {
		t.Errorf("Expected 0 services, got %d", len(tcpConfig.Services))
	}
	if len(tcpConfig.Middlewares) != 0 {
		t.Errorf("Expected 0 middlewares, got %d", len(tcpConfig.Middlewares))
	}

	// Test UDP config marshal errors
	udpConfig := &dynamic.UDPConfiguration{
		Routers:  make(map[string]*dynamic.UDPRouter),
		Services: make(map[string]*dynamic.UDPService),
	}

	udpProviderConfig := &config.UDPSection{
		Routers: &config.UDPRoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				// Channel cannot be marshaled - triggers marshal error continue
				make(chan int),
			},
		},
		Services: &config.UDPServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				// Channel cannot be marshaled - triggers marshal error continue
				make(chan int),
			},
		},
	}

	err = parseUDPConfig(map[string]interface{}{}, udpConfig, udpProviderConfig, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// All should be empty due to marshal errors
	if len(udpConfig.Routers) != 0 {
		t.Errorf("Expected 0 routers, got %d", len(udpConfig.Routers))
	}
	if len(udpConfig.Services) != 0 {
		t.Errorf("Expected 0 services, got %d", len(udpConfig.Services))
	}
}

func TestParseConfigUnmarshalErrors(t *testing.T) {
	// Test HTTP config unmarshal errors
	httpConfig := &dynamic.HTTPConfiguration{
		Routers:     make(map[string]*dynamic.Router),
		Services:    make(map[string]*dynamic.Service),
		Middlewares: make(map[string]*dynamic.Middleware),
	}

	httpProviderConfig := &config.HTTPSection{
		Routers: &config.RoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				map[string]interface{}{
					"name": "test-router",
					"rule": 123, // Invalid type - should cause unmarshal error
				},
			},
		},
		Services: &config.ServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				map[string]interface{}{
					"name":         "test-service",
					"loadBalancer": "invalid", // Invalid structure - should cause unmarshal error
				},
			},
		},
		Middlewares: &config.MiddlewaresConfig{
			Discover: true,
			ExtraMiddlewares: []interface{}{
				map[string]interface{}{
					"name":      "test-middleware",
					"basicAuth": "invalid", // Invalid structure - should cause unmarshal error
				},
			},
		},
	}

	err := parseHTTPConfig(map[string]interface{}{}, httpConfig, httpProviderConfig, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// All should be empty due to unmarshal errors
	if len(httpConfig.Routers) != 0 {
		t.Errorf("Expected 0 routers due to unmarshal error, got %d", len(httpConfig.Routers))
	}
	if len(httpConfig.Services) != 0 {
		t.Errorf("Expected 0 services due to unmarshal error, got %d", len(httpConfig.Services))
	}
	if len(httpConfig.Middlewares) != 0 {
		t.Errorf("Expected 0 middlewares due to unmarshal error, got %d", len(httpConfig.Middlewares))
	}

	// Test TCP config unmarshal errors
	tcpConfig := &dynamic.TCPConfiguration{
		Routers:     make(map[string]*dynamic.TCPRouter),
		Services:    make(map[string]*dynamic.TCPService),
		Middlewares: make(map[string]*dynamic.TCPMiddleware),
	}

	tcpProviderConfig := &config.TCPSection{
		Routers: &config.RoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				map[string]interface{}{
					"name": "test-tcp-router",
					"rule": 123, // Invalid type - should cause unmarshal error
				},
			},
		},
		Services: &config.ServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				map[string]interface{}{
					"name":         "test-tcp-service",
					"loadBalancer": "invalid", // Invalid structure - should cause unmarshal error
				},
			},
		},
		Middlewares: &config.MiddlewaresConfig{
			Discover: true,
			ExtraMiddlewares: []interface{}{
				map[string]interface{}{
					"name":        "test-tcp-middleware",
					"ipWhiteList": "invalid", // Invalid structure - should cause unmarshal error
				},
			},
		},
	}

	err = parseTCPConfig(map[string]interface{}{}, tcpConfig, tcpProviderConfig, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// All should be empty due to unmarshal errors
	if len(tcpConfig.Routers) != 0 {
		t.Errorf("Expected 0 routers due to unmarshal error, got %d", len(tcpConfig.Routers))
	}
	if len(tcpConfig.Services) != 0 {
		t.Errorf("Expected 0 services due to unmarshal error, got %d", len(tcpConfig.Services))
	}
	if len(tcpConfig.Middlewares) != 0 {
		t.Errorf("Expected 0 middlewares due to unmarshal error, got %d", len(tcpConfig.Middlewares))
	}

	// Test UDP config unmarshal errors
	udpConfig := &dynamic.UDPConfiguration{
		Routers:  make(map[string]*dynamic.UDPRouter),
		Services: make(map[string]*dynamic.UDPService),
	}

	udpProviderConfig := &config.UDPSection{
		Routers: &config.UDPRoutersConfig{
			Discover: true,
			ExtraRoutes: []interface{}{
				map[string]interface{}{
					"name":        "test-udp-router",
					"entryPoints": "invalid", // Invalid type - should cause unmarshal error
				},
			},
		},
		Services: &config.UDPServicesConfig{
			Discover: true,
			ExtraServices: []interface{}{
				map[string]interface{}{
					"name":         "test-udp-service",
					"loadBalancer": "invalid", // Invalid structure - should cause unmarshal error
				},
			},
		},
	}

	err = parseUDPConfig(map[string]interface{}{}, udpConfig, udpProviderConfig, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// All should be empty due to unmarshal errors
	if len(udpConfig.Routers) != 0 {
		t.Errorf("Expected 0 routers due to unmarshal error, got %d", len(udpConfig.Routers))
	}
	if len(udpConfig.Services) != 0 {
		t.Errorf("Expected 0 services due to unmarshal error, got %d", len(udpConfig.Services))
	}
}
