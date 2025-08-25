package filters

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func TestUDPRouters(t *testing.T) {
	tests := []struct {
		name     string
		routers  map[string]*dynamic.UDPRouter
		pattern  string
		expected []string
	}{
		{
			name: "filter all UDP routers",
			routers: map[string]*dynamic.UDPRouter{
				"udp-router": {
					Service: "udp-service",
				},
				"dns-router": {
					Service: "dns-service",
				},
				"admin": {
					Service: "admin",
				},
			},
			pattern:  ".*",
			expected: []string{"admin", "dns-router", "udp-router"},
		},
		{
			name: "filter specific pattern",
			routers: map[string]*dynamic.UDPRouter{
				"udp-router": {
					Service: "udp-service",
				},
				"dns-router": {
					Service: "dns-service",
				},
			},
			pattern:  "udp-.*",
			expected: []string{"udp-router"},
		},
		{
			name: "filter specific pattern",
			routers: map[string]*dynamic.UDPRouter{
				"udp-router": {
					Service: "udp-service",
				},
				"dns-router": {
					Service: "dns-service",
				},
			},
			pattern:  "dns-.*",
			expected: []string{"dns-router"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UDPRouters(tt.routers, &config.UDPRoutersConfig{Filter: config.UDPRouterFilter{Name: tt.pattern}}, config.ProviderFilter{})

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d routers, got %d", len(tt.expected), len(result))
				return
			}

			for _, expectedName := range tt.expected {
				if _, found := result[expectedName]; !found {
					t.Errorf("Expected router %s not found in result", expectedName)
				}
			}
		})
	}
}

func TestUDPServices(t *testing.T) {
	tests := []struct {
		name     string
		services map[string]*dynamic.UDPService
		pattern  string
		expected []string
	}{
		{
			name: "filter all UDP services",
			services: map[string]*dynamic.UDPService{
				"udp-service": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{Address: "udp1:53"},
						},
					},
				},
				"dns-service": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{Address: "dns1:53"},
						},
					},
				},
				"admin": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{Address: "admin:53"},
						},
					},
				},
			},
			pattern:  ".*",
			expected: []string{"admin", "dns-service", "udp-service"},
		},
		{
			name: "filter specific pattern",
			services: map[string]*dynamic.UDPService{
				"udp-service": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{Address: "udp1:53"},
						},
					},
				},
				"dns-service": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{Address: "dns1:53"},
						},
					},
				},
			},
			pattern:  "udp-.*",
			expected: []string{"udp-service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UDPServices(tt.services, &config.UDPServicesConfig{Filter: config.ServiceFilter{Name: tt.pattern}}, config.ProviderFilter{})

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d services, got %d", len(tt.expected), len(result))
				return
			}

			for _, expectedName := range tt.expected {
				if _, found := result[expectedName]; !found {
					t.Errorf("Expected service %s not found in result", expectedName)
				}
			}
		})
	}
}

func TestUDPRoutersAdvancedFiltering(t *testing.T) {
	t.Run("filter by entrypoints", func(t *testing.T) {
		routers := map[string]*dynamic.UDPRouter{
			"udp-router-1": {
				Service:     "udp-service-1",
				EntryPoints: []string{"udp", "udp-secure"},
			},
			"udp-router-2": {
				Service:     "udp-service-2",
				EntryPoints: []string{"udp-alt"},
			},
		}

		cfg := &config.UDPRoutersConfig{
			Filter: config.UDPRouterFilter{
				Entrypoints: []string{"udp"},
			},
		}

		result := UDPRouters(routers, cfg, config.ProviderFilter{})

		if len(result) != 1 {
			t.Errorf("Expected 1 router, got %d", len(result))
		}

		if _, exists := result["udp-router-1"]; !exists {
			t.Error("Expected udp-router-1 to be in result")
		}
	})

	t.Run("filter by service", func(t *testing.T) {
		routers := map[string]*dynamic.UDPRouter{
			"udp-router-1": {
				Service: "udp-service-1",
			},
			"udp-router-2": {
				Service: "udp-service-2",
			},
		}

		cfg := &config.UDPRoutersConfig{
			Filter: config.UDPRouterFilter{
				Service: "udp-service-1",
			},
		}

		result := UDPRouters(routers, cfg, config.ProviderFilter{})

		if len(result) != 1 {
			t.Errorf("Expected 1 router, got %d", len(result))
		}

		if _, exists := result["udp-router-1"]; !exists {
			t.Error("Expected udp-router-1 to be in result")
		}
	})

	t.Run("invalid service regex", func(t *testing.T) {
		routers := map[string]*dynamic.UDPRouter{
			"udp-router-1": {
				Service: "udp-service-1",
			},
		}

		cfg := &config.UDPRoutersConfig{
			Filter: config.UDPRouterFilter{
				Service: "*", // Invalid regex
			},
		}

		result := UDPRouters(routers, cfg, config.ProviderFilter{})

		if len(result) != 0 {
			t.Errorf("Expected 0 routers due to invalid regex, got %d", len(result))
		}
	})
}
