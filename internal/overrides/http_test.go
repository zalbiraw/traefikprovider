package overrides

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
)

func TestOverrideHTTPRouters_RuleOverride(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"test-router": {Rule: "Host(`example.com`)", Service: "test-service"},
	}
	overrides := config.RouterOverrides{
		Rules: []config.OverrideRule{{
			Matcher: "Name(`test-router`)",
			Value:   "Host(`new.example.com`)",
		}},
	}
	OverrideHTTPRouters(routers, overrides)
	if routers["test-router"].Rule != "Host(`new.example.com`)" {
		t.Errorf("Expected rule 'Host(`new.example.com`)', got %s", routers["test-router"].Rule)
	}
}

func TestOverrideHTTPRouters_RuleOverrideWithDollar(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"api-router": {Rule: "Host(`api.example.com`)", Service: "api-service"},
	}
	overrides := config.RouterOverrides{
		Rules: []config.OverrideRule{{
			Matcher: "Name(`api-router`)",
			Value:   "$1 && PathPrefix(`/v1`)",
		}},
	}
	OverrideHTTPRouters(routers, overrides)
	expected := "Host(`api.example.com`) && PathPrefix(`/v1`)"
	if routers["api-router"].Rule != expected {
		t.Errorf("Expected rule '%s', got %s", expected, routers["api-router"].Rule)
	}
}

func TestOverrideHTTPRouters_ServiceOverride(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"service-router": {Rule: "Host(`service.example.com`)", Service: "old-service"},
	}
	overrides := config.RouterOverrides{
		Services: []config.OverrideService{{
			Matcher: "Name(`service-router`)",
			Value:   "new-service",
		}},
	}
	OverrideHTTPRouters(routers, overrides)
	if routers["service-router"].Service != "new-service" {
		t.Errorf("Expected service 'new-service', got %s", routers["service-router"].Service)
	}
}

func TestOverrideHTTPRouters_ServiceOverrideWithDollar(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"service-router": {Rule: "Host(`service.example.com`)", Service: "old-service"},
	}
	overrides := config.RouterOverrides{
		Services: []config.OverrideService{{
			Matcher: "Name(`service-router`)",
			Value:   "$1-v2",
		}},
	}
	OverrideHTTPRouters(routers, overrides)
	if routers["service-router"].Service != "old-service-v2" {
		t.Errorf("Expected service 'old-service-v2', got %s", routers["service-router"].Service)
	}
}

func TestOverrideHTTPRouters_EntrypointsString(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"ep-router": {Rule: "Host(`ep.example.com`)", Service: "ep-service", EntryPoints: []string{"web"}},
	}
	overrides := config.RouterOverrides{
		Entrypoints: []config.OverrideEntrypoint{{
			Matcher: "Name(`ep-router`)",
			Value:   "websecure",
		}},
	}
	OverrideHTTPRouters(routers, overrides)
	expected := []string{"web", "websecure"}
	if len(routers["ep-router"].EntryPoints) != 2 || routers["ep-router"].EntryPoints[1] != "websecure" {
		t.Errorf("Expected entrypoints %v, got %v", expected, routers["ep-router"].EntryPoints)
	}
}

func TestOverrideHTTPRouters_EntrypointsArray(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"ep-router": {Rule: "Host(`ep.example.com`)", Service: "ep-service", EntryPoints: []string{"web"}},
	}
	overrides := config.RouterOverrides{
		Entrypoints: []config.OverrideEntrypoint{{
			Matcher: "Name(`ep-router`)",
			Value:   []string{"web", "websecure"},
		}},
	}
	OverrideHTTPRouters(routers, overrides)
	expected := []string{"web", "websecure"}
	if len(routers["ep-router"].EntryPoints) != 2 {
		t.Errorf("Expected entrypoints %v, got %v", expected, routers["ep-router"].EntryPoints)
	}
}

func TestOverrideHTTPRouters_MiddlewaresString(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"mw-router": {Rule: "Host(`mw.example.com`)", Service: "mw-service", Middlewares: []string{"auth"}},
	}
	overrides := config.RouterOverrides{
		Middlewares: []config.OverrideMiddleware{{
			Matcher: "Name(`mw-router`)",
			Value:   "cors",
		}},
	}
	OverrideHTTPRouters(routers, overrides)
	expected := []string{"auth", "cors"}
	if len(routers["mw-router"].Middlewares) != 2 || routers["mw-router"].Middlewares[1] != "cors" {
		t.Errorf("Expected middlewares %v, got %v", expected, routers["mw-router"].Middlewares)
	}
}

func TestOverrideHTTPRouters_MiddlewaresArray(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"mw-router": {Rule: "Host(`mw.example.com`)", Service: "mw-service", Middlewares: []string{"auth"}},
	}
	overrides := config.RouterOverrides{
		Middlewares: []config.OverrideMiddleware{{
			Matcher: "Name(`mw-router`)",
			Value:   []string{"cors", "ratelimit"},
		}},
	}
	OverrideHTTPRouters(routers, overrides)
	expected := []string{"cors", "ratelimit"}
	if len(routers["mw-router"].Middlewares) != 2 || routers["mw-router"].Middlewares[0] != "cors" || routers["mw-router"].Middlewares[1] != "ratelimit" {
		t.Errorf("Expected middlewares %v, got %v", expected, routers["mw-router"].Middlewares)
	}
}

func TestOverrideHTTPServices_ServerOverrideString(t *testing.T) {
	services := map[string]*dynamic.Service{
		"test-service": {
			LoadBalancer: &dynamic.ServersLoadBalancer{
				Servers: []dynamic.Server{{URL: "http://old-server:8080"}},
			},
		},
	}

	overrides := config.ServiceOverrides{
		Servers: []config.OverrideServer{{
			Matcher: "Name(`test-service`)",
			Value:   "http://new-server:8080",
		}},
	}

	OverrideHTTPServices(services, overrides)

	if len(services["test-service"].LoadBalancer.Servers) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(services["test-service"].LoadBalancer.Servers))
	}
	if services["test-service"].LoadBalancer.Servers[1].URL != "http://new-server:8080" {
		t.Errorf("Expected second server 'http://new-server:8080', got %s", services["test-service"].LoadBalancer.Servers[1].URL)
	}
}

func TestOverrideHTTPServices_WithTunnel(t *testing.T) {
	services := map[string]*dynamic.Service{
		"tunnel-service": {
			LoadBalancer: &dynamic.ServersLoadBalancer{
				Servers: []dynamic.Server{{URL: "http://old-server:8080"}},
			},
		},
	}

	tunnels := []config.TunnelConfig{{
		Name:      "my-tunnel",
		Addresses: []string{"http://tunnel1:8080", "http://tunnel2:8080"},
	}}

	overrides := config.ServiceOverrides{
		Servers: []config.OverrideServer{{
			Matcher: "Name(`tunnel-service`)",
			Value:   []string{"http://ignored:8080"},
			Tunnel:  "my-tunnel",
		}},
	}

	OverrideHTTPServices(services, overrides)

	// Apply tunnels in a separate pass per new design
	ApplyTunnels(services, overrides, tunnels)

	// Tunnels should be used for HTTP services when specified.
	if len(services["tunnel-service"].LoadBalancer.Servers) != 2 {
		t.Errorf("Expected 2 servers from tunnel, got %d", len(services["tunnel-service"].LoadBalancer.Servers))
	}
	if services["tunnel-service"].LoadBalancer.Servers[0].URL != "http://tunnel1:8080" {
		t.Errorf("Expected first server 'http://tunnel1:8080', got %s", services["tunnel-service"].LoadBalancer.Servers[0].URL)
	}
}

func TestOverrideHTTPServices_HealthcheckOverride(t *testing.T) {
	services := map[string]*dynamic.Service{
		"health-service": {
			LoadBalancer: &dynamic.ServersLoadBalancer{
				Servers:     []dynamic.Server{{URL: "http://server:8080"}},
				HealthCheck: &dynamic.ServerHealthCheck{Path: "/old-health", Interval: "10s", Timeout: "2s"},
			},
		},
	}

	overrides := config.ServiceOverrides{
		Healthchecks: []config.OverrideHealthcheck{{
			Matcher:  "Name(`health-service`)",
			Path:     "/health",
			Interval: "30s",
			Timeout:  "5s",
		}},
	}

	OverrideHTTPServices(services, overrides)

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
	if hc.Timeout != "5s" {
		t.Errorf("Expected health check timeout '5s', got %s", hc.Timeout)
	}
}

func TestOverrideHTTPServices_HealthcheckPartial(t *testing.T) {
	services := map[string]*dynamic.Service{
		"health-service": {
			LoadBalancer: &dynamic.ServersLoadBalancer{
				Servers:     []dynamic.Server{{URL: "http://server:8080"}},
				HealthCheck: &dynamic.ServerHealthCheck{Path: "/old-health", Interval: "10s", Timeout: "2s"},
			},
		},
	}

	overrides := config.ServiceOverrides{
		Healthchecks: []config.OverrideHealthcheck{{
			Matcher: "Name(`health-service`)",
			Path:    "/new-health", // Only override path
		}},
	}

	OverrideHTTPServices(services, overrides)

	hc := services["health-service"].LoadBalancer.HealthCheck
	if hc.Path != "/new-health" {
		t.Errorf("Expected health check path '/new-health', got %s", hc.Path)
	}
	if hc.Interval != "10s" {
		t.Errorf("Expected health check interval '10s' (unchanged), got %s", hc.Interval)
	}
	if hc.Timeout != "2s" {
		t.Errorf("Expected health check timeout '2s' (unchanged), got %s", hc.Timeout)
	}
}

func TestOverrideHTTPServices_NoHealthcheck(t *testing.T) {
	services := map[string]*dynamic.Service{
		"no-hc-service": {
			LoadBalancer: &dynamic.ServersLoadBalancer{
				Servers: []dynamic.Server{{URL: "http://server:8080"}},
			},
		},
	}

	overrides := config.ServiceOverrides{
		Healthchecks: []config.OverrideHealthcheck{{
			Matcher: "Name(`no-hc-service`)",
			Path:    "/health",
		}},
	}

	// Should not panic when health check is nil
	OverrideHTTPServices(services, overrides)

	// Service should remain unchanged
	if services["no-hc-service"].LoadBalancer.HealthCheck != nil {
		t.Error("Expected health check to remain nil")
	}
}
