package filters

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func TestFilterMapByNameRegex(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]*struct{}
		pattern  string
		expected []string
	}{
		{
			name: "match all services",
			input: map[string]*struct{}{
				"web-app":     {},
				"api-service": {},
				"database":    {},
			},
			pattern:  ".*",
			expected: []string{"api-service", "database", "web-app"},
		},
		{
			name: "match specific pattern",
			input: map[string]*struct{}{
				"web-app":     {},
				"api-service": {},
				"database":    {},
			},
			pattern:  "web-.*",
			expected: []string{"web-app"},
		},
		{
			name: "include all services",
			input: map[string]*struct{}{
				"web-service":    {},
				"api-service":    {},
				"admin@internal": {},
				"debug@internal": {},
			},
			pattern:  ".*",
			expected: []string{"admin@internal", "api-service", "debug@internal", "web-service"},
		},
		{
			name: "no matches",
			input: map[string]*struct{}{
				"web-app":     {},
				"api-service": {},
			},
			pattern:  "nonexistent",
			expected: []string{},
		},
		{
			name:     "empty input",
			input:    map[string]*struct{}{},
			pattern:  ".*",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterMapByNameRegex[struct{}, *struct{}](tt.input, tt.pattern, "")

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d results, got %d", len(tt.expected), len(result))
				return
			}

			// Convert result to slice for comparison
			resultSlice := make([]string, 0, len(result))
			for key := range result {
				resultSlice = append(resultSlice, key)
			}

			// Sort both slices for comparison
			for _, expected := range tt.expected {
				found := false
				for _, actual := range resultSlice {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected key %s not found in result", expected)
				}
			}
		})
	}
}

func TestFilterMapByNameRegexInvalidPattern(t *testing.T) {
	input := map[string]*struct{}{
		"service1": {},
	}

	// Invalid regex pattern should return empty result
	result := filterMapByNameRegex[struct{}, *struct{}](input, "[", "")
	if len(result) != 0 {
		t.Errorf("Expected empty result for invalid regex, got %d items", len(result))
	}
}

func TestRegexMatchInvalidPattern(t *testing.T) {
	// Test invalid regex pattern
	matched, err := regexMatch("[", "test")
	if err == nil {
		t.Error("Expected error for invalid regex pattern")
	}
	if matched {
		t.Error("Expected no match for invalid regex")
	}
}

func TestRegexMatchEmptyPattern(t *testing.T) {
	// Test empty pattern should return true
	matched, err := regexMatch("", "test")
	if err != nil {
		t.Errorf("Expected no error for empty pattern, got: %v", err)
	}
	if !matched {
		t.Error("Expected match for empty pattern")
	}
}

func TestRouterEntrypointsMatch(t *testing.T) {
	tests := []struct {
		name      string
		routerEPs []string
		filterEPs []string
		expected  bool
	}{
		{
			name:      "empty filter entrypoints",
			routerEPs: []string{"web", "websecure"},
			filterEPs: []string{},
			expected:  true,
		},
		{
			name:      "router contains all filter entrypoints",
			routerEPs: []string{"web", "websecure", "api"},
			filterEPs: []string{"web", "websecure"},
			expected:  true,
		},
		{
			name:      "router missing some filter entrypoints",
			routerEPs: []string{"web"},
			filterEPs: []string{"web", "websecure"},
			expected:  false,
		},
		{
			name:      "empty router entrypoints",
			routerEPs: []string{},
			filterEPs: []string{"web"},
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := routerEntrypointsMatch(tt.routerEPs, tt.filterEPs)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHTTPRoutersWithEntrypoints(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"web-router": {
			Rule:        "Host(`example.com`)",
			EntryPoints: []string{"web"},
			Service:     "my-service",
		},
		"api-router": {
			Rule:        "PathPrefix(`/api`)",
			EntryPoints: []string{"api"},
			Service:     "api-service",
		},
	}

	// Test entrypoint filtering
	result := HTTPRouters(routers, &config.RoutersConfig{
		Filter: config.RouterFilter{
			Name:        ".*",
			Entrypoints: []string{"web"},
		},
	}, config.ProviderFilter{})

	if len(result) != 1 {
		t.Errorf("Expected 1 router with web entrypoint, got %d", len(result))
	}

	if _, found := result["web-router"]; !found {
		t.Error("Expected web-router to be found")
	}
}

func TestHTTPRoutersWithRuleFilter(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"host-router": {
			Rule:    "Host(`example.com`)",
			Service: "host-service",
		},
		"path-router": {
			Rule:    "PathPrefix(`/api`)",
			Service: "path-service",
		},
	}

	// Test rule filtering
	result := HTTPRouters(routers, &config.RoutersConfig{
		Filter: config.RouterFilter{
			Name: ".*",
			Rule: "Host\\(.*\\)",
		},
	}, config.ProviderFilter{})

	if len(result) != 1 {
		t.Errorf("Expected 1 router with Host rule, got %d", len(result))
	}

	if _, found := result["host-router"]; !found {
		t.Error("Expected host-router to be found")
	}
}

func TestHTTPRoutersWithServiceFilter(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"service1-router": {
			Rule:    "Host(`example.com`)",
			Service: "my-service",
		},
		"service2-router": {
			Rule:    "PathPrefix(`/api`)",
			Service: "api-service",
		},
	}

	// Test service filtering
	result := HTTPRouters(routers, &config.RoutersConfig{
		Filter: config.RouterFilter{
			Name:    ".*",
			Service: "^my-.*",
		},
	}, config.ProviderFilter{})

	if len(result) != 1 {
		t.Errorf("Expected 1 router with my-service, got %d", len(result))
	}

	if _, found := result["service1-router"]; !found {
		t.Error("Expected service1-router to be found")
	}
}

func TestHTTPRoutersWithInvalidData(t *testing.T) {
	// Test with valid router data
	routers := map[string]*dynamic.Router{
		"valid-router": {
			Rule:    "Host(`example.com`)",
			Service: "my-service",
		},
	}

	result := HTTPRouters(routers, &config.RoutersConfig{
		Filter: config.RouterFilter{Name: ".*"},
	}, config.ProviderFilter{})

	if len(result) != 1 {
		t.Errorf("Expected 1 valid router, got %d", len(result))
	}

	if _, found := result["valid-router"]; !found {
		t.Error("Expected valid-router to be found")
	}
}

func TestTCPRoutersWithEntrypoints(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"tcp-web": {
			Rule:        "HostSNI(`example.com`)",
			EntryPoints: []string{"tcp-web", "tcp-secure"},
			Service:     "tcp-service",
		},
		"tcp-api": {
			Rule:        "HostSNI(`api.example.com`)",
			EntryPoints: []string{"tcp-api"},
			Service:     "tcp-api-service",
		},
	}

	result := TCPRouters(routers, &config.RoutersConfig{
		Filter: config.RouterFilter{
			Name:        ".*",
			Entrypoints: []string{"tcp-web"},
		},
	}, config.ProviderFilter{})

	if len(result) != 1 {
		t.Errorf("Expected 1 TCP router with tcp-web entrypoint, got %d", len(result))
	}

	if _, found := result["tcp-web"]; !found {
		t.Error("Expected tcp-web router to be found")
	}
}

func TestTCPRoutersWithRuleFilter(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"sni-router": {
			Rule:    "HostSNI(`example.com`)",
			Service: "sni-service",
		},
		"catch-all": {
			Rule:    "HostSNI(`*`)",
			Service: "catch-service",
		},
	}

	result := TCPRouters(routers, &config.RoutersConfig{
		Filter: config.RouterFilter{
			Name: ".*",
			Rule: "HostSNI\\(`\\*`\\)",
		},
	}, config.ProviderFilter{})

	if len(result) != 1 {
		t.Errorf("Expected 1 TCP router with catch-all rule, got %d", len(result))
	}

	if _, found := result["catch-all"]; !found {
		t.Error("Expected catch-all router to be found")
	}
}

func TestUDPRoutersWithServiceFilter(t *testing.T) {
	routers := map[string]*dynamic.UDPRouter{
		"dns-router": {
			Service: "dns-service",
		},
		"syslog-router": {
			Service: "syslog-service",
		},
	}

	result := UDPRouters(routers, &config.UDPRoutersConfig{
		Filter: config.UDPRouterFilter{
			Name:    ".*",
			Service: "^dns-.*",
		},
	}, config.ProviderFilter{})

	if len(result) != 1 {
		t.Errorf("Expected 1 UDP router with dns service, got %d", len(result))
	}

	if _, found := result["dns-router"]; !found {
		t.Error("Expected dns-router to be found")
	}
}

func TestInvalidInputHandling(t *testing.T) {
	// Test with empty maps - the functions now expect typed maps, not interface{}

	// HTTP functions
	httpResult := HTTPRouters(map[string]*dynamic.Router{}, &config.RoutersConfig{}, config.ProviderFilter{})
	if len(httpResult) != 0 {
		t.Error("Expected empty result for empty HTTP routers input")
	}

	httpServicesResult := HTTPServices(map[string]*dynamic.Service{}, &config.ServicesConfig{}, config.ProviderFilter{})
	if len(httpServicesResult) != 0 {
		t.Error("Expected empty result for empty HTTP services input")
	}

	httpMiddlewaresResult := HTTPMiddlewares(map[string]*dynamic.Middleware{}, &config.MiddlewaresConfig{}, config.ProviderFilter{})
	if len(httpMiddlewaresResult) != 0 {
		t.Error("Expected empty result for empty HTTP middlewares input")
	}

	// TCP functions
	tcpResult := TCPRouters(map[string]*dynamic.TCPRouter{}, &config.RoutersConfig{}, config.ProviderFilter{})
	if len(tcpResult) != 0 {
		t.Error("Expected empty result for empty TCP routers input")
	}

	tcpServicesResult := TCPServices(map[string]*dynamic.TCPService{}, &config.ServicesConfig{}, config.ProviderFilter{})
	if len(tcpServicesResult) != 0 {
		t.Error("Expected empty result for empty TCP services input")
	}

	tcpMiddlewaresResult := TCPMiddlewares(map[string]*dynamic.TCPMiddleware{}, &config.MiddlewaresConfig{}, config.ProviderFilter{})
	if len(tcpMiddlewaresResult) != 0 {
		t.Error("Expected empty result for empty TCP middlewares input")
	}

	// UDP functions
	udpResult := UDPRouters(map[string]*dynamic.UDPRouter{}, &config.UDPRoutersConfig{}, config.ProviderFilter{})
	if len(udpResult) != 0 {
		t.Error("Expected empty result for empty UDP routers input")
	}

	udpServicesResult := UDPServices(map[string]*dynamic.UDPService{}, &config.UDPServicesConfig{}, config.ProviderFilter{})
	if len(udpServicesResult) != 0 {
		t.Error("Expected empty result for empty UDP services input")
	}
}

func TestExtractProviderFromName(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"noat", ""},
		{"svc@file", "file"},
		{"ns/name@kubernetes@file", "file"},
		{"trailing@", ""},
	}
	for _, c := range cases {
		if got := extractProviderFromName(c.in); got != c.out {
			t.Fatalf("extractProviderFromName(%q)=%q want %q", c.in, got, c.out)
		}
	}
}

func TestFilterMapByNameRegexProviderInvalid(t *testing.T) {
	input := map[string]*struct{}{
		"a@p": {},
	}
	// invalid provider regex should yield empty result
	res := filterMapByNameRegex[struct{}, *struct{}](input, "", "[")
	if len(res) != 0 {
		t.Fatalf("expected 0, got %d", len(res))
	}
}

func TestHTTPRoutersProviderOverride(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"r1@p1": {Rule: "Host(`a`)", Service: "s1"},
		"r2@p2": {Rule: "Host(`b`)", Service: "s2"},
	}
	cfg := &config.RoutersConfig{Filter: config.RouterFilter{Name: ".*", Provider: "p1"}}
	// ProviderFilter should override filter.Provider to p2
	out := HTTPRouters(routers, cfg, config.ProviderFilter{Provider: "p2"})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if _, ok := out["r2@p2"]; !ok {
		t.Fatalf("expected r2@p2 to remain after provider override")
	}
}

func TestHTTPRoutersDiscoverPriorityFalse(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"r@p": {Rule: "Host(`x`)", Service: "s", Priority: 5, EntryPoints: []string{"web"}},
	}
	cfg := &config.RoutersConfig{DiscoverPriority: false, Filter: config.RouterFilter{Name: ".*", Entrypoints: []string{"web"}}}
	out := HTTPRouters(routers, cfg, config.ProviderFilter{})
	r := out["r@p"]
	if r == nil || r.Priority != 0 {
		t.Fatalf("expected priority reset to 0, got %+v", r)
	}
}

func TestHTTPServicesEarlyReturn(t *testing.T) {
	services := map[string]*dynamic.Service{
		"a": {},
		"b": {},
	}
	// Empty filter -> early return full map
	out := HTTPServices(services, &config.ServicesConfig{Filter: config.ServiceFilter{}}, config.ProviderFilter{})
	if len(out) != 2 {
		t.Fatalf("expected 2 services, got %d", len(out))
	}
}

func TestHTTPMiddlewaresEarlyReturn(t *testing.T) {
	m := map[string]*dynamic.Middleware{
		"m1": {},
		"m2": {},
	}
	out := HTTPMiddlewares(m, &config.MiddlewaresConfig{Filter: config.MiddlewareFilter{}}, config.ProviderFilter{})
	if len(out) != 2 {
		t.Fatalf("expected 2 middlewares, got %d", len(out))
	}
}

func TestHTTPServicesProviderOverride(t *testing.T) {
	services := map[string]*dynamic.Service{
		"s1@p1": {},
		"s2@p2": {},
	}
	cfg := &config.ServicesConfig{Filter: config.ServiceFilter{Name: ".*", Provider: "p1"}}
	out := HTTPServices(services, cfg, config.ProviderFilter{Provider: "p2"})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if _, ok := out["s2@p2"]; !ok {
		t.Fatalf("expected s2@p2 to remain after provider override")
	}
}

func TestHTTPMiddlewaresProviderOverride(t *testing.T) {
	m := map[string]*dynamic.Middleware{
		"m1@p1": {},
		"m2@p2": {},
	}
	cfg := &config.MiddlewaresConfig{Filter: config.MiddlewareFilter{Name: ".*", Provider: "p1"}}
	out := HTTPMiddlewares(m, cfg, config.ProviderFilter{Provider: "p2"})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if _, ok := out["m2@p2"]; !ok {
		t.Fatalf("expected m2@p2 to remain after provider override")
	}
}

func TestTCPServicesProviderOverride(t *testing.T) {
	services := map[string]*dynamic.TCPService{
		"ts1@p1": {},
		"ts2@p2": {},
	}
	cfg := &config.ServicesConfig{Filter: config.ServiceFilter{Name: ".*", Provider: "p1"}}
	out := TCPServices(services, cfg, config.ProviderFilter{Provider: "p2"})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if _, ok := out["ts2@p2"]; !ok {
		t.Fatalf("expected ts2@p2 to remain after provider override")
	}
}

func TestTCPMiddlewaresProviderOverride(t *testing.T) {
	m := map[string]*dynamic.TCPMiddleware{
		"tm1@p1": {},
		"tm2@p2": {},
	}
	cfg := &config.MiddlewaresConfig{Filter: config.MiddlewareFilter{Name: ".*", Provider: "p1"}}
	out := TCPMiddlewares(m, cfg, config.ProviderFilter{Provider: "p2"})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if _, ok := out["tm2@p2"]; !ok {
		t.Fatalf("expected tm2@p2 to remain after provider override")
	}
}

func TestTCPRoutersProviderOverride(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"tr1@p1": {Rule: "HostSNI(`*`)", Service: "s1"},
		"tr2@p2": {Rule: "HostSNI(`*`)", Service: "s2"},
	}
	cfg := &config.RoutersConfig{Filter: config.RouterFilter{Name: ".*", Provider: "p1"}}
	out := TCPRouters(routers, cfg, config.ProviderFilter{Provider: "p2"})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if _, ok := out["tr2@p2"]; !ok {
		t.Fatalf("expected tr2@p2 to remain after provider override")
	}
}

func TestUDPRoutersProviderOverride(t *testing.T) {
	routers := map[string]*dynamic.UDPRouter{
		"ur1@p1": {Service: "s1"},
		"ur2@p2": {Service: "s2"},
	}
	cfg := &config.UDPRoutersConfig{Filter: config.UDPRouterFilter{Name: ".*", Provider: "p1"}}
	out := UDPRouters(routers, cfg, config.ProviderFilter{Provider: "p2"})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if _, ok := out["ur2@p2"]; !ok {
		t.Fatalf("expected ur2@p2 to remain after provider override")
	}
}

func TestUDPServicesProviderOverride(t *testing.T) {
	services := map[string]*dynamic.UDPService{
		"us1@p1": {},
		"us2@p2": {},
	}
	cfg := &config.UDPServicesConfig{Filter: config.ServiceFilter{Name: ".*", Provider: "p1"}}
	out := UDPServices(services, cfg, config.ProviderFilter{Provider: "p2"})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if _, ok := out["us2@p2"]; !ok {
		t.Fatalf("expected us2@p2 to remain after provider override")
	}
}
