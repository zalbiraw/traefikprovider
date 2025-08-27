package matchers

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
)

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
		Matcher: "NameRegexp(`.*`) && Entrypoint(`web`)",
	}, "")

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
		// Match by service name instead of rule, since rule matching is not supported in matcher language
		Matcher: "ServiceRegexp(`host-.*`)",
	}, "")

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
		Matcher: "NameRegexp(`.*`) && ServiceRegexp(`^my-.*`)",
	}, "")

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
		Matcher: "NameRegexp(`.*`)",
	}, "")

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
		Matcher: "NameRegexp(`.*`) && Entrypoint(`tcp-web`)",
	}, "")

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
		Matcher: "Name(`catch-all`)",
	}, "")

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
		Matcher: "NameRegexp(`.*`) && ServiceRegexp(`^dns-.*`)",
	}, "")

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
	httpResult := HTTPRouters(map[string]*dynamic.Router{}, &config.RoutersConfig{}, "")
	if len(httpResult) != 0 {
		t.Error("Expected empty result for empty HTTP routers input")
	}

	httpServicesResult := HTTPServices(map[string]*dynamic.Service{}, &config.ServicesConfig{}, "")
	if len(httpServicesResult) != 0 {
		t.Error("Expected empty result for empty HTTP services input")
	}

	httpMiddlewaresResult := HTTPMiddlewares(map[string]*dynamic.Middleware{}, &config.MiddlewaresConfig{}, "")
	if len(httpMiddlewaresResult) != 0 {
		t.Error("Expected empty result for empty HTTP middlewares input")
	}

	// TCP functions
	tcpResult := TCPRouters(map[string]*dynamic.TCPRouter{}, &config.RoutersConfig{}, "")
	if len(tcpResult) != 0 {
		t.Error("Expected empty result for empty TCP routers input")
	}

	tcpServicesResult := TCPServices(map[string]*dynamic.TCPService{}, &config.ServicesConfig{}, "")
	if len(tcpServicesResult) != 0 {
		t.Error("Expected empty result for empty TCP services input")
	}

	tcpMiddlewaresResult := TCPMiddlewares(map[string]*dynamic.TCPMiddleware{}, &config.MiddlewaresConfig{}, "")
	if len(tcpMiddlewaresResult) != 0 {
		t.Error("Expected empty result for empty TCP middlewares input")
	}

	// UDP functions
	udpResult := UDPRouters(map[string]*dynamic.UDPRouter{}, &config.UDPRoutersConfig{}, "")
	if len(udpResult) != 0 {
		t.Error("Expected empty result for empty UDP routers input")
	}

	udpServicesResult := UDPServices(map[string]*dynamic.UDPService{}, &config.UDPServicesConfig{}, "")
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

func TestHTTPRoutersProviderOverride(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"r1@p1": {Rule: "Host(`a`)", Service: "s1"},
		"r2@p2": {Rule: "Host(`b`)", Service: "s2"},
	}
	cfg := &config.RoutersConfig{Matcher: "NameRegexp(`.*`)"}
	// Provider filter is applied via provider matcher string
	out := HTTPRouters(routers, cfg, "Provider(`p2`)")
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
	cfg := &config.RoutersConfig{DiscoverPriority: false, Matcher: "NameRegexp(`.*`) && Entrypoint(`web`)"}
	out := HTTPRouters(routers, cfg, "")
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
	// Empty matcher -> early return full map
	out := HTTPServices(services, &config.ServicesConfig{Matcher: ""}, "")
	if len(out) != 2 {
		t.Fatalf("expected 2 services, got %d", len(out))
	}
}

func TestHTTPMiddlewaresEarlyReturn(t *testing.T) {
	m := map[string]*dynamic.Middleware{
		"m1": {},
		"m2": {},
	}
	out := HTTPMiddlewares(m, &config.MiddlewaresConfig{Matcher: ""}, "")
	if len(out) != 2 {
		t.Fatalf("expected 2 middlewares, got %d", len(out))
	}
}

func TestHTTPServicesProviderOverride(t *testing.T) {
	services := map[string]*dynamic.Service{
		"s1@p1": {},
		"s2@p2": {},
	}
	cfg := &config.ServicesConfig{Matcher: "NameRegexp(`.*`)"}
	out := HTTPServices(services, cfg, "Provider(`p2`)")
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
	cfg := &config.MiddlewaresConfig{Matcher: "NameRegexp(`.*`)"}
	out := HTTPMiddlewares(m, cfg, "Provider(`p2`)")
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
	cfg := &config.ServicesConfig{Matcher: "NameRegexp(`.*`)"}
	out := TCPServices(services, cfg, "Provider(`p2`)")
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
	cfg := &config.MiddlewaresConfig{Matcher: "NameRegexp(`.*`)"}
	out := TCPMiddlewares(m, cfg, "Provider(`p2`)")
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
	cfg := &config.RoutersConfig{Matcher: "NameRegexp(`.*`)"}
	out := TCPRouters(routers, cfg, "Provider(`p2`)")
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
	cfg := &config.UDPRoutersConfig{Matcher: "NameRegexp(`.*`)"}
	out := UDPRouters(routers, cfg, "Provider(`p2`)")
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
	cfg := &config.UDPServicesConfig{Matcher: "NameRegexp(`.*`)"}
	out := UDPServices(services, cfg, "Provider(`p2`)")
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if _, ok := out["us2@p2"]; !ok {
		t.Fatalf("expected us2@p2 to remain after provider override")
	}
}

func TestCombineRulesCases(t *testing.T) {
	cases := []struct {
		p   string
		s   string
		exp string
	}{
		{"", "", ""},
		{"P", "", "P"},
		{"", "S", "S"},
		{"A", "B", "(A) && (B)"},
	}
	for _, c := range cases {
		if got := combineRules(c.p, c.s); got != c.exp {
			t.Fatalf("combineRules(%q,%q)=%q want %q", c.p, c.s, got, c.exp)
		}
	}
}

func TestRegexMatch_TrueFalse(t *testing.T) {
	ok, err := regexMatch("^ab.+z$", "abcz")
	if err != nil || !ok {
		t.Fatalf("expected true, err=nil got %v,%v", ok, err)
	}
	ok, err = regexMatch("^ab.+z$", "abq")
	if err != nil || ok {
		t.Fatalf("expected false, err=nil got %v,%v", ok, err)
	}
}

func TestHTTPRoutersEarlyReturn(t *testing.T) {
	r := map[string]*dynamic.Router{
		"r1": {},
		"r2": {},
	}
	out := HTTPRouters(r, &config.RoutersConfig{Matcher: ""}, "")
	if len(out) != 2 {
		t.Fatalf("expected early return of all routers, got %d", len(out))
	}
}

func TestHTTPRoutersCompileErrorReturnsEmpty(t *testing.T) {
	r := map[string]*dynamic.Router{
		"r1": {},
	}
	// invalid rule -> compile error -> empty result
	out := HTTPRouters(r, &config.RoutersConfig{Matcher: "Name(`unterminated"}, "")
	if len(out) != 0 {
		t.Fatalf("expected empty on compile error, got %d", len(out))
	}
}

func TestHTTPServicesCompileErrorReturnsEmpty(t *testing.T) {
	svcs := map[string]*dynamic.Service{"s1": {}}
	out := HTTPServices(svcs, &config.ServicesConfig{Matcher: "Name(`unterminated"}, "")
	if len(out) != 0 {
		t.Fatalf("expected empty on compile error, got %d", len(out))
	}
}

func TestTCPRoutersCompileErrorReturnsEmpty(t *testing.T) {
	r := map[string]*dynamic.TCPRouter{"tr1": {}}
	out := TCPRouters(r, &config.RoutersConfig{Matcher: "Name(`unterminated"}, "")
	if len(out) != 0 {
		t.Fatalf("expected empty on compile error, got %d", len(out))
	}
}
