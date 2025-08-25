package internal

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	tlstypes "github.com/traefik/genconf/dynamic/tls"
)

func TestMergeConfigurations(t *testing.T) {
	tests := []struct {
		name     string
		configs  []*dynamic.Configuration
		expected *dynamic.Configuration
	}{
		{
			name:    "merge empty configurations",
			configs: []*dynamic.Configuration{},
			expected: &dynamic.Configuration{
				HTTP: &dynamic.HTTPConfiguration{
					Routers:           map[string]*dynamic.Router{},
					Services:          map[string]*dynamic.Service{},
					Middlewares:       map[string]*dynamic.Middleware{},
					ServersTransports: map[string]*dynamic.ServersTransport{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:     map[string]*dynamic.TCPRouter{},
					Services:    map[string]*dynamic.TCPService{},
					Middlewares: map[string]*dynamic.TCPMiddleware{},
				},
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Certificates: []*tlstypes.CertAndStores{},
					Options:      map[string]tlstypes.Options{},
					Stores:       map[string]tlstypes.Store{},
				},
			},
		},
		{
			name: "merge single configuration",
			configs: []*dynamic.Configuration{
				{
					HTTP: &dynamic.HTTPConfiguration{
						Routers: map[string]*dynamic.Router{
							"router1": {
								Rule:    "Host(`example.com`)",
								Service: "service1",
							},
						},
						Services: map[string]*dynamic.Service{
							"service1": {
								LoadBalancer: &dynamic.ServersLoadBalancer{
									Servers: []dynamic.Server{
										{URL: "http://backend:8080"},
									},
								},
							},
						},
					},
				},
			},
			expected: &dynamic.Configuration{
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"router1": {
							Rule:    "Host(`example.com`)",
							Service: "service1",
						},
					},
					Services: map[string]*dynamic.Service{
						"service1": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{URL: "http://backend:8080"},
								},
							},
						},
					},
					Middlewares:       map[string]*dynamic.Middleware{},
					ServersTransports: map[string]*dynamic.ServersTransport{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:     map[string]*dynamic.TCPRouter{},
					Services:    map[string]*dynamic.TCPService{},
					Middlewares: map[string]*dynamic.TCPMiddleware{},
				},
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Certificates: []*tlstypes.CertAndStores{},
					Options:      map[string]tlstypes.Options{},
					Stores:       map[string]tlstypes.Store{},
				},
			},
		},
		{
			name: "merge multiple configurations with HTTP",
			configs: []*dynamic.Configuration{
				{
					HTTP: &dynamic.HTTPConfiguration{
						Routers: map[string]*dynamic.Router{
							"router1": {
								Rule:    "Host(`api.example.com`)",
								Service: "api-service",
							},
						},
						Services: map[string]*dynamic.Service{
							"api-service": {
								LoadBalancer: &dynamic.ServersLoadBalancer{
									Servers: []dynamic.Server{
										{URL: "http://api-backend:8080"},
									},
								},
							},
						},
					},
				},
				{
					HTTP: &dynamic.HTTPConfiguration{
						Routers: map[string]*dynamic.Router{
							"router2": {
								Rule:    "Host(`web.example.com`)",
								Service: "web-service",
							},
						},
						Services: map[string]*dynamic.Service{
							"web-service": {
								LoadBalancer: &dynamic.ServersLoadBalancer{
									Servers: []dynamic.Server{
										{URL: "http://web-backend:8080"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "merge configurations with TCP",
			configs: []*dynamic.Configuration{
				{
					TCP: &dynamic.TCPConfiguration{
						Routers: map[string]*dynamic.TCPRouter{
							"tcp-router1": {
								Rule:    "HostSNI(`tcp.example.com`)",
								Service: "tcp-service1",
							},
						},
						Services: map[string]*dynamic.TCPService{
							"tcp-service1": {
								LoadBalancer: &dynamic.TCPServersLoadBalancer{
									Servers: []dynamic.TCPServer{
										{Address: "tcp-backend:8081"},
									},
								},
							},
						},
					},
				},
				{
					TCP: &dynamic.TCPConfiguration{
						Routers: map[string]*dynamic.TCPRouter{
							"tcp-router2": {
								Rule:    "HostSNI(`tcp2.example.com`)",
								Service: "tcp-service2",
							},
						},
						Services: map[string]*dynamic.TCPService{
							"tcp-service2": {
								LoadBalancer: &dynamic.TCPServersLoadBalancer{
									Servers: []dynamic.TCPServer{
										{Address: "tcp2-backend:8081"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "merge configurations with UDP",
			configs: []*dynamic.Configuration{
				{
					UDP: &dynamic.UDPConfiguration{
						Routers: map[string]*dynamic.UDPRouter{
							"udp-router1": {
								Service: "udp-service1",
							},
						},
						Services: map[string]*dynamic.UDPService{
							"udp-service1": {
								LoadBalancer: &dynamic.UDPServersLoadBalancer{
									Servers: []dynamic.UDPServer{
										{Address: "udp-backend:8082"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "merge configurations with TLS",
			configs: []*dynamic.Configuration{
				{
					TLS: &dynamic.TLSConfiguration{
						Certificates: []*tlstypes.CertAndStores{
							{
								Certificate: tlstypes.Certificate{
									CertFile: "/path/to/cert1.pem",
									KeyFile:  "/path/to/key1.pem",
								},
							},
						},
						Options: map[string]tlstypes.Options{
							"default": {
								MinVersion: "VersionTLS12",
							},
						},
						Stores: map[string]tlstypes.Store{
							"default": {
								DefaultCertificate: &tlstypes.Certificate{
									CertFile: "/path/to/default.pem",
									KeyFile:  "/path/to/default-key.pem",
								},
							},
						},
					},
				},
				{
					TLS: &dynamic.TLSConfiguration{
						Certificates: []*tlstypes.CertAndStores{
							{
								Certificate: tlstypes.Certificate{
									CertFile: "/path/to/cert2.pem",
									KeyFile:  "/path/to/key2.pem",
								},
							},
						},
						Options: map[string]tlstypes.Options{
							"custom": {
								MinVersion: "VersionTLS13",
							},
						},
					},
				},
			},
		},
		{
			name: "merge configurations with nil sections",
			configs: []*dynamic.Configuration{
				{
					HTTP: &dynamic.HTTPConfiguration{
						Routers: map[string]*dynamic.Router{
							"router1": {
								Rule:    "Host(`example.com`)",
								Service: "service1",
							},
						},
					},
					TCP: nil,
					UDP: nil,
					TLS: nil,
				},
				{
					HTTP: nil,
					TCP: &dynamic.TCPConfiguration{
						Routers: map[string]*dynamic.TCPRouter{
							"tcp-router1": {
								Rule:    "HostSNI(`tcp.example.com`)",
								Service: "tcp-service1",
							},
						},
					},
					UDP: nil,
					TLS: nil,
				},
			},
		},
		{
			name: "merge with nil configurations",
			configs: []*dynamic.Configuration{
				{
					HTTP: &dynamic.HTTPConfiguration{
						Routers: map[string]*dynamic.Router{
							"router1": {
								Rule:    "Host(`example.com`)",
								Service: "service1",
							},
						},
					},
				},
				nil,
				{
					HTTP: &dynamic.HTTPConfiguration{
						Routers: map[string]*dynamic.Router{
							"router2": {
								Rule:    "Host(`example2.com`)",
								Service: "service2",
							},
						},
					},
				},
			},
		},
		{
			name: "merge configurations with overlapping keys",
			configs: []*dynamic.Configuration{
				{
					HTTP: &dynamic.HTTPConfiguration{
						Routers: map[string]*dynamic.Router{
							"router1": {
								Rule:    "Host(`example.com`)",
								Service: "service1",
							},
						},
						Services: map[string]*dynamic.Service{
							"shared-service": {
								LoadBalancer: &dynamic.ServersLoadBalancer{
									Servers: []dynamic.Server{
										{URL: "http://backend1:8080"},
									},
								},
							},
						},
					},
				},
				{
					HTTP: &dynamic.HTTPConfiguration{
						Routers: map[string]*dynamic.Router{
							"router1": { // Same key, should overwrite
								Rule:    "Host(`updated.example.com`)",
								Service: "updated-service",
							},
						},
						Services: map[string]*dynamic.Service{
							"shared-service": { // Same key, should overwrite
								LoadBalancer: &dynamic.ServersLoadBalancer{
									Servers: []dynamic.Server{
										{URL: "http://backend2:8080"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "merge all protocol types",
			configs: []*dynamic.Configuration{
				{
					HTTP: &dynamic.HTTPConfiguration{
						Routers: map[string]*dynamic.Router{
							"http-router": {
								Rule:    "Host(`http.example.com`)",
								Service: "http-service",
							},
						},
						Services: map[string]*dynamic.Service{
							"http-service": {
								LoadBalancer: &dynamic.ServersLoadBalancer{
									Servers: []dynamic.Server{
										{URL: "http://http-backend:8080"},
									},
								},
							},
						},
						Middlewares: map[string]*dynamic.Middleware{
							"auth": {
								BasicAuth: &dynamic.BasicAuth{
									Users: []string{"user:password"},
								},
							},
						},
						ServersTransports: map[string]*dynamic.ServersTransport{
							"default": {
								ServerName: "example.com",
							},
						},
					},
					TCP: &dynamic.TCPConfiguration{
						Routers: map[string]*dynamic.TCPRouter{
							"tcp-router": {
								Rule:    "HostSNI(`tcp.example.com`)",
								Service: "tcp-service",
							},
						},
						Services: map[string]*dynamic.TCPService{
							"tcp-service": {
								LoadBalancer: &dynamic.TCPServersLoadBalancer{
									Servers: []dynamic.TCPServer{
										{Address: "tcp-backend:8081"},
									},
								},
							},
						},
						Middlewares: map[string]*dynamic.TCPMiddleware{
							"tcp-auth": {
								IPWhiteList: &dynamic.TCPIPWhiteList{
									SourceRange: []string{"192.168.1.0/24"},
								},
							},
						},
					},
					UDP: &dynamic.UDPConfiguration{
						Routers: map[string]*dynamic.UDPRouter{
							"udp-router": {
								Service: "udp-service",
							},
						},
						Services: map[string]*dynamic.UDPService{
							"udp-service": {
								LoadBalancer: &dynamic.UDPServersLoadBalancer{
									Servers: []dynamic.UDPServer{
										{Address: "udp-backend:8082"},
									},
								},
							},
						},
					},
					TLS: &dynamic.TLSConfiguration{
						Certificates: []*tlstypes.CertAndStores{
							{
								Certificate: tlstypes.Certificate{
									CertFile: "/path/to/cert.pem",
									KeyFile:  "/path/to/key.pem",
								},
							},
						},
						Options: map[string]tlstypes.Options{
							"default": {
								MinVersion: "VersionTLS12",
							},
						},
						Stores: map[string]tlstypes.Store{
							"default": {
								DefaultCertificate: &tlstypes.Certificate{
									CertFile: "/path/to/default.pem",
									KeyFile:  "/path/to/default-key.pem",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeConfigurations(tt.configs...)

			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			// Verify structure is always initialized
			if result.HTTP == nil {
				t.Error("Expected HTTP configuration to be initialized")
			}
			if result.TCP == nil {
				t.Error("Expected TCP configuration to be initialized")
			}
			if result.UDP == nil {
				t.Error("Expected UDP configuration to be initialized")
			}
			if result.TLS == nil {
				t.Error("Expected TLS configuration to be initialized")
			}

			// For specific test cases, verify expected content
			if tt.expected != nil {
				if tt.expected.HTTP != nil {
					if len(tt.expected.HTTP.Routers) > 0 && len(result.HTTP.Routers) != len(tt.expected.HTTP.Routers) {
						t.Errorf("Expected %d HTTP routers, got %d", len(tt.expected.HTTP.Routers), len(result.HTTP.Routers))
					}
					if len(tt.expected.HTTP.Services) > 0 && len(result.HTTP.Services) != len(tt.expected.HTTP.Services) {
						t.Errorf("Expected %d HTTP services, got %d", len(tt.expected.HTTP.Services), len(result.HTTP.Services))
					}
				}
			}

			// Verify maps are always initialized (never nil)
			if result.HTTP.Routers == nil {
				t.Error("Expected HTTP routers map to be initialized")
			}
			if result.HTTP.Services == nil {
				t.Error("Expected HTTP services map to be initialized")
			}
			if result.HTTP.Middlewares == nil {
				t.Error("Expected HTTP middlewares map to be initialized")
			}
			if result.HTTP.ServersTransports == nil {
				t.Error("Expected HTTP server transports map to be initialized")
			}
			if result.TCP.Routers == nil {
				t.Error("Expected TCP routers map to be initialized")
			}
			if result.TCP.Services == nil {
				t.Error("Expected TCP services map to be initialized")
			}
			if result.TCP.Middlewares == nil {
				t.Error("Expected TCP middlewares map to be initialized")
			}
			if result.UDP.Routers == nil {
				t.Error("Expected UDP routers map to be initialized")
			}
			if result.UDP.Services == nil {
				t.Error("Expected UDP services map to be initialized")
			}
			if result.TLS.Certificates == nil {
				t.Error("Expected TLS certificates slice to be initialized")
			}
			if result.TLS.Options == nil {
				t.Error("Expected TLS options map to be initialized")
			}
			if result.TLS.Stores == nil {
				t.Error("Expected TLS stores map to be initialized")
			}
		})
	}
}

func TestMergeConfigurationsOverwrite(t *testing.T) {
	// Test that later configurations overwrite earlier ones for the same keys
	config1 := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"shared-router": {
					Rule:    "Host(`old.example.com`)",
					Service: "old-service",
				},
			},
			Services: map[string]*dynamic.Service{
				"shared-service": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{
							{URL: "http://old-backend:8080"},
						},
					},
				},
			},
		},
	}

	config2 := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"shared-router": {
					Rule:    "Host(`new.example.com`)",
					Service: "new-service",
				},
			},
			Services: map[string]*dynamic.Service{
				"shared-service": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{
							{URL: "http://new-backend:8080"},
						},
					},
				},
			},
		},
	}

	result := MergeConfigurations(config1, config2)

	// Verify the second configuration overwrote the first
	if result.HTTP.Routers["shared-router"].Rule != "Host(`new.example.com`)" {
		t.Error("Expected router to be overwritten by later configuration")
	}
	if result.HTTP.Routers["shared-router"].Service != "new-service" {
		t.Error("Expected router service to be overwritten by later configuration")
	}

	if result.HTTP.Services["shared-service"].LoadBalancer.Servers[0].URL != "http://new-backend:8080" {
		t.Error("Expected service to be overwritten by later configuration")
	}
}

func TestMergeConfigurationsTLSCertificateAppend(t *testing.T) {
	// Test that TLS certificates are appended, not overwritten
	config1 := &dynamic.Configuration{
		TLS: &dynamic.TLSConfiguration{
			Certificates: []*tlstypes.CertAndStores{
				{
					Certificate: tlstypes.Certificate{
						CertFile: "/path/to/cert1.pem",
						KeyFile:  "/path/to/key1.pem",
					},
				},
			},
		},
	}

	config2 := &dynamic.Configuration{
		TLS: &dynamic.TLSConfiguration{
			Certificates: []*tlstypes.CertAndStores{
				{
					Certificate: tlstypes.Certificate{
						CertFile: "/path/to/cert2.pem",
						KeyFile:  "/path/to/key2.pem",
					},
				},
			},
		},
	}

	result := MergeConfigurations(config1, config2)

	// Verify certificates are appended
	if len(result.TLS.Certificates) != 2 {
		t.Errorf("Expected 2 certificates, got %d", len(result.TLS.Certificates))
	}

	if result.TLS.Certificates[0].Certificate.CertFile != "/path/to/cert1.pem" {
		t.Error("Expected first certificate to be preserved")
	}
	if result.TLS.Certificates[1].Certificate.CertFile != "/path/to/cert2.pem" {
		t.Error("Expected second certificate to be appended")
	}
}
