package filters

import (
	"testing"

	"github.com/zalbiraw/traefik-provider/config"
)

func TestFilterHTTPRouters(t *testing.T) {
	tests := []struct {
		name     string
		routers  map[string]interface{}
		pattern  string
		expected []string
	}{
		{
			name: "filter all routers",
			routers: map[string]interface{}{
				"web-router":     map[string]interface{}{"rule": "Host(`web.example.com`)", "service": "web-service"},
				"api-router":     map[string]interface{}{"rule": "Host(`api.example.com`)", "service": "api-service"},
				"admin@internal": map[string]interface{}{"rule": "Host(`admin.traefik`)", "service": "admin@internal"},
			},
			pattern:  ".*",
			expected: []string{"api-router", "web-router"},
		},
		{
			name: "filter specific pattern",
			routers: map[string]interface{}{
				"web-router": map[string]interface{}{"rule": "Host(`web.example.com`)", "service": "web-service"},
				"api-router": map[string]interface{}{"rule": "Host(`api.example.com`)", "service": "api-service"},
			},
			pattern:  "web-.*",
			expected: []string{"web-router"},
		},
		{
			name: "exclude internal routers",
			routers: map[string]interface{}{
				"web-router":     map[string]interface{}{"rule": "Host(`web.example.com`)", "service": "web-service"},
				"admin@internal": map[string]interface{}{"rule": "Host(`admin.traefik`)", "service": "admin@internal"},
			},
			pattern:  ".*",
			expected: []string{"web-router"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTTPRouters(tt.routers, &config.RoutersConfig{Filters: config.RouterFilters{Name: tt.pattern}})

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

func TestFilterHTTPServices(t *testing.T) {
	tests := []struct {
		name     string
		services map[string]interface{}
		pattern  string
		expected []string
	}{
		{
			name: "filter all services",
			services: map[string]interface{}{
				"web-service":     map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"url": "http://web:80"}}}},
				"api-service":     map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"url": "http://api:80"}}}},
				"admin@internal":  map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"url": "http://admin:80"}}}},
			},
			pattern:  ".*",
			expected: []string{"api-service", "web-service"},
		},
		{
			name: "filter specific pattern",
			services: map[string]interface{}{
				"web-service": map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"url": "http://web:80"}}}},
				"api-service": map[string]interface{}{"loadBalancer": map[string]interface{}{"servers": []interface{}{map[string]interface{}{"url": "http://api:80"}}}},
			},
			pattern:  "web-.*",
			expected: []string{"web-service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTTPServices(tt.services, &config.ServicesConfig{Filters: config.ServiceFilters{Name: tt.pattern}})

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

func TestFilterHTTPMiddlewares(t *testing.T) {
	tests := []struct {
		name        string
		middlewares map[string]interface{}
		pattern     string
		expected    []string
	}{
		{
			name: "filter all middlewares",
			middlewares: map[string]interface{}{
				"auth-middleware": map[string]interface{}{"basicAuth": map[string]interface{}{"users": []string{"user:pass"}}},
				"cors-middleware": map[string]interface{}{"headers": map[string]interface{}{"accessControlAllowOriginList": []string{"*"}}},
			},
			pattern:  ".*",
			expected: []string{"auth-middleware", "cors-middleware"},
		},
		{
			name: "filter specific pattern",
			middlewares: map[string]interface{}{
				"auth-middleware": map[string]interface{}{"basicAuth": map[string]interface{}{"users": []string{"user:pass"}}},
				"cors-middleware": map[string]interface{}{"headers": map[string]interface{}{"accessControlAllowOriginList": []string{"*"}}},
			},
			pattern:  "auth-.*",
			expected: []string{"auth-middleware"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTTPMiddlewares(tt.middlewares, &config.MiddlewaresConfig{Filters: config.MiddlewareFilters{Name: tt.pattern}})

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

func TestFilterHTTPServersTransports(t *testing.T) {
	tests := []struct {
		name       string
		transports map[string]interface{}
		pattern    string
		expected   []string
	}{
		{
			name: "filter all transports",
			transports: map[string]interface{}{
				"secure-transport":   map[string]interface{}{"serverName": "secure.example.com"},
				"insecure-transport": map[string]interface{}{"insecureSkipVerify": true},
			},
			pattern:  ".*",
			expected: []string{"insecure-transport", "secure-transport"},
		},
		{
			name: "filter specific pattern",
			transports: map[string]interface{}{
				"secure-transport":   map[string]interface{}{"serverName": "secure.example.com"},
				"insecure-transport": map[string]interface{}{"insecureSkipVerify": true},
			},
			pattern:  "^secure-transport$",
			expected: []string{"secure-transport"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTTPServerTransports(tt.transports, &config.ServerTransportsConfig{Filters: config.ServerTransportFilters{Name: tt.pattern}})

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d transports, got %d", len(tt.expected), len(result))
				return
			}

			for _, expectedName := range tt.expected {
				if _, found := result[expectedName]; !found {
					t.Errorf("Expected transport %s not found in result", expectedName)
				}
			}
		})
	}
}
