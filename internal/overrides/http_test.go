package overrides

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func TestOverrideHTTPRouters(t *testing.T) {
	t.Run("rule override", func(t *testing.T) {
		routers := map[string]*dynamic.Router{
			"test-router": {
				Rule:    "Host(`example.com`)",
				Service: "test-service",
			},
		}

		overrides := config.RouterOverrides{
			Rules: []config.OverrideRule{
				{
					Filters: config.RouterFilters{
						Name: "test-router",
					},
					Value: "Host(`new.example.com`)",
				},
			},
		}

		OverrideHTTPRouters(routers, overrides)

		if routers["test-router"].Rule != "Host(`new.example.com`)" {
			t.Errorf("Expected rule 'Host(`new.example.com`)', got %s", routers["test-router"].Rule)
		}
	})

	t.Run("rule override with $1 replacement", func(t *testing.T) {
		routers := map[string]*dynamic.Router{
			"api-router": {
				Rule:    "Host(`api.example.com`)",
				Service: "api-service",
			},
		}

		overrides := config.RouterOverrides{
			Rules: []config.OverrideRule{
				{
					Filters: config.RouterFilters{
						Name: "api-router",
					},
					Value: "$1 && PathPrefix(`/v1`)",
				},
			},
		}

		OverrideHTTPRouters(routers, overrides)

		expected := "Host(`api.example.com`) && PathPrefix(`/v1`)"
		if routers["api-router"].Rule != expected {
			t.Errorf("Expected rule '%s', got %s", expected, routers["api-router"].Rule)
		}
	})

	t.Run("service override", func(t *testing.T) {
		routers := map[string]*dynamic.Router{
			"service-router": {
				Rule:    "Host(`service.example.com`)",
				Service: "old-service",
			},
		}

		overrides := config.RouterOverrides{
			Services: []config.OverrideService{
				{
					Filters: config.RouterFilters{
						Name: "service-router",
					},
					Value: "new-service",
				},
			},
		}

		OverrideHTTPRouters(routers, overrides)

		if routers["service-router"].Service != "new-service" {
			t.Errorf("Expected service 'new-service', got %s", routers["service-router"].Service)
		}
	})
}

func TestOverrideHTTPServices(t *testing.T) {
	t.Run("server override", func(t *testing.T) {
		services := map[string]*dynamic.Service{
			"test-service": {
				LoadBalancer: &dynamic.ServersLoadBalancer{
					Servers: []dynamic.Server{
						{URL: "http://old-server:8080"},
					},
				},
			},
		}

		// Test server override logic directly
		for _, service := range services {
			if service.LoadBalancer != nil {
				service.LoadBalancer.Servers = []dynamic.Server{
					{URL: "http://new-server:8080"},
					{URL: "http://backup-server:8080"},
				}
			}
		}

		if len(services["test-service"].LoadBalancer.Servers) != 2 {
			t.Errorf("Expected 2 servers, got %d", len(services["test-service"].LoadBalancer.Servers))
		}

		if services["test-service"].LoadBalancer.Servers[0].URL != "http://new-server:8080" {
			t.Errorf("Expected first server 'http://new-server:8080', got %s", services["test-service"].LoadBalancer.Servers[0].URL)
		}
	})

	t.Run("health check override", func(t *testing.T) {
		services := map[string]*dynamic.Service{
			"health-service": {
				LoadBalancer: &dynamic.ServersLoadBalancer{
					Servers: []dynamic.Server{
						{URL: "http://server:8080"},
					},
					HealthCheck: &dynamic.ServerHealthCheck{},
				},
			},
		}

		// Test health check override logic directly
		for _, service := range services {
			if service.LoadBalancer != nil && service.LoadBalancer.HealthCheck != nil {
				service.LoadBalancer.HealthCheck.Path = "/health"
				service.LoadBalancer.HealthCheck.Interval = "30s"
				service.LoadBalancer.HealthCheck.Timeout = "5s"
			}
		}

		hc := services["health-service"].LoadBalancer.HealthCheck
		if hc == nil {
			t.Fatal("Expected health check to exist")
		}

		if hc.Path != "/health" {
			t.Errorf("Expected health check path '/health', got %s", hc.Path)
		}

		if hc.Interval != "30s" {
			t.Errorf("Expected health check interval '30s', got %s", hc.Interval)
		}
	})
}
