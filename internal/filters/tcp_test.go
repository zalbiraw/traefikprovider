package filters

import (
	"testing"

	"github.com/zalbiraw/traefik-provider/config"
)

func TestTCPRouters(t *testing.T) {
	tests := []struct {
		name     string
		routers  map[string]interface{}
		pattern  string
		expected []string
	}{
		{
			name: "filter all TCP routers",
			routers: map[string]interface{}{
				"tcp-router":     map[string]interface{}{"rule": "HostSNI(`tcp.example.com`)", "service": "tcp-service"},
				"mysql-router":   map[string]interface{}{"rule": "HostSNI(`mysql.example.com`)", "service": "mysql-service"},
				"admin@internal": map[string]interface{}{"rule": "HostSNI(`admin.traefik`)", "service": "admin@internal"},
			},
			pattern:  ".*",
			expected: []string{"mysql-router", "tcp-router"},
		},
		{
			name: "filter specific pattern",
			routers: map[string]interface{}{
				"tcp-router":   map[string]interface{}{"rule": "HostSNI(`tcp.example.com`)", "service": "tcp-service"},
				"mysql-router": map[string]interface{}{"rule": "HostSNI(`mysql.example.com`)", "service": "mysql-service"},
			},
			pattern:  "tcp-.*",
			expected: []string{"tcp-router"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TCPRouters(tt.routers, &config.RoutersConfig{Filters: config.RouterFilters{Name: tt.pattern}})

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

func TestTCPServices(t *testing.T) {
	tests := []struct {
		name     string
		services map[string]interface{}
		pattern  string
		expected []string
	}{
		{
			name: "filter all TCP services",
			services: map[string]interface{}{
				"tcp-service":     map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"address": "tcp1:3306"}}}},
				"mysql-service":   map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"address": "mysql1:3306"}}}},
				"admin@internal":  map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"address": "admin:3306"}}}},
			},
			pattern:  ".*",
			expected: []string{"mysql-service", "tcp-service"},
		},
		{
			name: "filter specific pattern",
			services: map[string]interface{}{
				"tcp-service":   map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"address": "tcp1:3306"}}}},
				"mysql-service": map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"address": "mysql1:3306"}}}},
			},
			pattern:  "tcp-.*",
			expected: []string{"tcp-service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TCPServices(tt.services, &config.ServicesConfig{Filters: config.ServiceFilters{Name: tt.pattern}})

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

func TestTCPMiddlewares(t *testing.T) {
	tests := []struct {
		name        string
		middlewares map[string]interface{}
		pattern     string
		expected    []string
	}{
		{
			name: "filter all TCP middlewares",
			middlewares: map[string]interface{}{
				"tcp-auth":     map[string]interface{}{"ipWhiteList": map[string]interface{}{"sourceRange": []string{"192.168.1.0/24"}}},
				"tcp-limiter":  map[string]interface{}{"inFlightConn": map[string]interface{}{"amount": 100}},
			},
			pattern:  ".*",
			expected: []string{"tcp-auth", "tcp-limiter"},
		},
		{
			name: "filter specific pattern",
			middlewares: map[string]interface{}{
				"tcp-auth":    map[string]interface{}{"ipWhiteList": map[string]interface{}{"sourceRange": []string{"192.168.1.0/24"}}},
				"tcp-limiter": map[string]interface{}{"inFlightConn": map[string]interface{}{"amount": 100}},
			},
			pattern:  "^tcp-auth$",
			expected: []string{"tcp-auth"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TCPMiddlewares(tt.middlewares, &config.MiddlewaresConfig{Filters: config.MiddlewareFilters{Name: tt.pattern}})

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d middlewares, got %d", len(tt.expected), len(result))
				return
			}

			for _, expectedName := range tt.expected {
				if _, found := result[expectedName]; !found {
					t.Errorf("Expected middleware %s not found in result", expectedName)
				}
			}
		})
	}
}
