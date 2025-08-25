package filters

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func TestFilterTCPRouters(t *testing.T) {
	tests := []struct {
		name     string
		routers  map[string]*dynamic.TCPRouter
		pattern  string
		expected []string
	}{
		{
			name: "filter all routers",
			routers: map[string]*dynamic.TCPRouter{
				"tcp-router": {
					Rule:    "HostSNI(`tcp.example.com`)",
					Service: "tcp-service",
				},
				"mysql-router": {
					Rule:    "HostSNI(`mysql.example.com`)",
					Service: "mysql-service",
				},
			},
			pattern:  ".*",
			expected: []string{"mysql-router", "tcp-router"},
		},
		{
			name: "filter specific pattern",
			routers: map[string]*dynamic.TCPRouter{
				"tcp-router": {
					Rule:    "HostSNI(`tcp.example.com`)",
					Service: "tcp-service",
				},
				"mysql-router": {
					Rule:    "HostSNI(`mysql.example.com`)",
					Service: "mysql-service",
				},
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
		services map[string]*dynamic.TCPService
		pattern  string
		expected []string
	}{
		{
			name: "filter all TCP services",
			services: map[string]*dynamic.TCPService{
				"tcp-service": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{Address: "tcp1:3306"},
						},
					},
				},
				"mysql-service": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{Address: "mysql1:3306"},
						},
					},
				},
			},
			pattern:  ".*",
			expected: []string{"mysql-service", "tcp-service"},
		},
		{
			name: "filter specific pattern",
			services: map[string]*dynamic.TCPService{
				"tcp-service": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{Address: "tcp1:3306"},
						},
					},
				},
				"mysql-service": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{Address: "mysql1:3306"},
						},
					},
				},
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
		middlewares map[string]*dynamic.TCPMiddleware
		pattern     string
		expected    []string
	}{
		{
			name: "filter all TCP middlewares",
			middlewares: map[string]*dynamic.TCPMiddleware{
				"tcp-auth": {
					IPWhiteList: &dynamic.TCPIPWhiteList{
						SourceRange: []string{"192.168.1.0/24"},
					},
				},
				"tcp-limiter": {
					InFlightConn: &dynamic.TCPInFlightConn{
						Amount: 100,
					},
				},
			},
			pattern:  ".*",
			expected: []string{"tcp-auth", "tcp-limiter"},
		},
		{
			name: "filter specific pattern",
			middlewares: map[string]*dynamic.TCPMiddleware{
				"tcp-auth": {
					IPWhiteList: &dynamic.TCPIPWhiteList{
						SourceRange: []string{"192.168.1.0/24"},
					},
				},
				"tcp-limiter": {
					InFlightConn: &dynamic.TCPInFlightConn{
						Amount: 100,
					},
				},
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
