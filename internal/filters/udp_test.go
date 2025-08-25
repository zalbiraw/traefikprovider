package filters

import (
	"testing"

	"github.com/zalbiraw/traefik-provider/config"
)

func TestUDPRouters(t *testing.T) {
	tests := []struct {
		name     string
		routers  map[string]interface{}
		pattern  string
		expected []string
	}{
		{
			name: "filter all UDP routers",
			routers: map[string]interface{}{
				"udp-router":     map[string]interface{}{"service": "udp-service"},
				"dns-router":     map[string]interface{}{"service": "dns-service"},
				"admin@internal": map[string]interface{}{"service": "admin@internal"},
			},
			pattern:  ".*",
			expected: []string{"dns-router", "udp-router"},
		},
		{
			name: "filter specific pattern",
			routers: map[string]interface{}{
				"udp-router": map[string]interface{}{"service": "udp-service"},
				"dns-router": map[string]interface{}{"service": "dns-service"},
			},
			pattern:  "udp-.*",
			expected: []string{"udp-router"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UDPRouters(tt.routers, &config.UDPRoutersConfig{Filters: config.UDPRouterFilters{Name: tt.pattern}})

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
		services map[string]interface{}
		pattern  string
		expected []string
	}{
		{
			name: "filter all UDP services",
			services: map[string]interface{}{
				"udp-service":     map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"address": "udp1:53"}}}},
				"dns-service":     map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"address": "dns1:53"}}}},
				"admin@internal":  map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"address": "admin:53"}}}},
			},
			pattern:  ".*",
			expected: []string{"dns-service", "udp-service"},
		},
		{
			name: "filter specific pattern",
			services: map[string]interface{}{
				"udp-service": map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"address": "udp1:53"}}}},
				"dns-service": map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"address": "dns1:53"}}}},
			},
			pattern:  "udp-.*",
			expected: []string{"udp-service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UDPServices(tt.services, &config.UDPServicesConfig{Filters: config.ServiceFilters{Name: tt.pattern}})

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
