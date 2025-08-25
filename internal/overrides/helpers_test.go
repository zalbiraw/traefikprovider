package overrides

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func TestApplyRouterOverride(t *testing.T) {
	filtered := map[string]*dynamic.Router{
		"test-router": {
			Rule:    "Host(`example.com`)",
			Service: "test-service",
		},
	}

	filter := config.RouterFilter{}

	value := "new-rule"

	applyRouterOverride(filtered, filter, value, func(r *dynamic.Router, v string) {
		r.Rule = v
	})

	// The function should work through the filter system
	// This is more of an integration test
	if len(filtered) == 0 {
		t.Error("Expected router to remain in filtered map")
	}
}

func TestHandleRouterOverride(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{
			name:  "string value",
			value: "web",
		},
		{
			name:  "array value",
			value: []string{"web", "websecure"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := map[string]*dynamic.Router{
				"test-router": {
					Rule:        "Host(`example.com`)",
					Service:     "test-service",
					EntryPoints: []string{},
				},
			}

			filter := config.RouterFilter{}

			handleRouterOverride(filtered, filter, tt.value,
				func(r *dynamic.Router, arr []string) { r.EntryPoints = arr },
				func(r *dynamic.Router, s string) { r.EntryPoints = []string{s} },
			)

			// Test that function executes without error
			if len(filtered) == 0 {
				t.Error("Expected router to remain in filtered map")
			}
		})
	}
}

func TestApplyServiceOverride(t *testing.T) {
	filtered := map[string]*dynamic.Service{
		"test-service": {
			LoadBalancer: &dynamic.ServersLoadBalancer{
				Servers: []dynamic.Server{
					{URL: "http://old-server:8080"},
				},
			},
		},
	}

	filter := config.ServiceFilter{}

	value := []string{"http://new-server:8080"}

	applyServiceOverride(filtered, filter, value, func(s *dynamic.Service, urls []string) {
		if s.LoadBalancer != nil {
			s.LoadBalancer.Servers = make([]dynamic.Server, len(urls))
			for i, url := range urls {
				s.LoadBalancer.Servers[i] = dynamic.Server{URL: url}
			}
		}
	})

	// Test that function executes without error
	if len(filtered) == 0 {
		t.Error("Expected service to remain in filtered map")
	}
}

func TestApplyTCPServiceOverride(t *testing.T) {
	filtered := map[string]*dynamic.TCPService{
		"tcp-service": {
			LoadBalancer: &dynamic.TCPServersLoadBalancer{
				Servers: []dynamic.TCPServer{
					{Address: "old-server:8080"},
				},
			},
		},
	}

	filter := config.ServiceFilter{}

	value := []string{"new-server:8080"}

	applyTCPServiceOverride(filtered, filter, value, func(s *dynamic.TCPService, addresses []string) {
		if s.LoadBalancer != nil {
			s.LoadBalancer.Servers = make([]dynamic.TCPServer, len(addresses))
			for i, addr := range addresses {
				s.LoadBalancer.Servers[i] = dynamic.TCPServer{Address: addr}
			}
		}
	})

	// Test that function executes without error
	if len(filtered) == 0 {
		t.Error("Expected TCP service to remain in filtered map")
	}
}

func TestApplyUDPServiceOverride(t *testing.T) {
	filtered := map[string]*dynamic.UDPService{
		"udp-service": {
			LoadBalancer: &dynamic.UDPServersLoadBalancer{
				Servers: []dynamic.UDPServer{
					{Address: "old-server:8080"},
				},
			},
		},
	}

	filter := config.ServiceFilter{}

	value := []string{"new-server:8080"}

	applyUDPServiceOverride(filtered, filter, value, func(s *dynamic.UDPService, addresses []string) {
		if s.LoadBalancer != nil {
			s.LoadBalancer.Servers = make([]dynamic.UDPServer, len(addresses))
			for i, addr := range addresses {
				s.LoadBalancer.Servers[i] = dynamic.UDPServer{Address: addr}
			}
		}
	})

	// Test that function executes without error
	if len(filtered) == 0 {
		t.Error("Expected UDP service to remain in filtered map")
	}
}
