package filters

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func TestFilterHTTPRouters(t *testing.T) {
	tests := []struct {
		name     string
		routers  map[string]*dynamic.Router
		pattern  string
		expected []string
	}{
		{
			name: "filter all routers",
			routers: map[string]*dynamic.Router{
				"web-router": {Rule: "Host(`web.example.com`)", Service: "web-service"},
				"api-router": {Rule: "Host(`api.example.com`)", Service: "api-service"},
			},
			pattern:  ".*",
			expected: []string{"api-router", "web-router"},
		},
		{
			name: "filter specific pattern",
			routers: map[string]*dynamic.Router{
				"web-router": {Rule: "Host(`web.example.com`)", Service: "web-service"},
				"api-router": {Rule: "Host(`api.example.com`)", Service: "api-service"},
			},
			pattern:  "web-.*",
			expected: []string{"web-router"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTTPRouters(tt.routers, &config.RoutersConfig{Filter: config.RouterFilter{Name: tt.pattern}}, config.ProviderFilter{})

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
		services map[string]*dynamic.Service
		pattern  string
		expected []string
	}{
		{
			name: "filter all services",
			services: map[string]*dynamic.Service{
				"web-service": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{{URL: "http://web:80"}},
					},
				},
				"api-service": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{{URL: "http://api:80"}},
					},
				},
			},
			pattern:  ".*",
			expected: []string{"api-service", "web-service"},
		},
		{
			name: "filter specific pattern",
			services: map[string]*dynamic.Service{
				"web-service": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{{URL: "http://web:80"}},
					},
				},
				"api-service": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{{URL: "http://api:80"}},
					},
				},
			},
			pattern:  "web-.*",
			expected: []string{"web-service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTTPServices(tt.services, &config.ServicesConfig{Filter: config.ServiceFilter{Name: tt.pattern}}, config.ProviderFilter{})

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
		middlewares map[string]*dynamic.Middleware
		pattern     string
		expected    []string
	}{
		{
			name: "filter all middlewares",
			middlewares: map[string]*dynamic.Middleware{
				"auth-middleware": {
					BasicAuth: &dynamic.BasicAuth{
						Users: []string{"user:pass"},
					},
				},
				"cors-middleware": {
					Headers: &dynamic.Headers{
						AccessControlAllowOriginList: []string{"*"},
					},
				},
			},
			pattern:  ".*",
			expected: []string{"auth-middleware", "cors-middleware"},
		},
		{
			name: "filter specific pattern",
			middlewares: map[string]*dynamic.Middleware{
				"auth-middleware": {
					BasicAuth: &dynamic.BasicAuth{
						Users: []string{"user:pass"},
					},
				},
				"cors-middleware": {
					Headers: &dynamic.Headers{
						AccessControlAllowOriginList: []string{"*"},
					},
				},
			},
			pattern:  "auth-.*",
			expected: []string{"auth-middleware"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTTPMiddlewares(tt.middlewares, &config.MiddlewaresConfig{Filter: config.MiddlewareFilter{Name: tt.pattern}}, config.ProviderFilter{})

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
