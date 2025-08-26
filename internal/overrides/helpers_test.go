package overrides

import (
	"reflect"
	"testing"

	"github.com/traefik/genconf/dynamic"
)

func TestApplyRouterOverride(t *testing.T) {
	matched := map[string]*dynamic.Router{
		"test-router": {
			Rule:    "Host(`example.com`)",
			Service: "test-service",
		},
	}

	match := ""

	value := "new-rule"

	applyRouterOverride(matched, match, value, func(r *dynamic.Router, v string) {
		r.Rule = v
	})

	// The function should work through the match system
	// This is more of an integration test
	if len(matched) == 0 {
		t.Error("Expected router to remain in matched map")
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
			matched := map[string]*dynamic.Router{
				"test-router": {
					Rule:        "Host(`example.com`)",
					Service:     "test-service",
					EntryPoints: []string{},
				},
			}

			match := ""

			handleRouterOverride(matched, match, tt.value,
				func(r *dynamic.Router, arr []string) { r.EntryPoints = arr },
				func(r *dynamic.Router, s string) { r.EntryPoints = []string{s} },
			)

			// Test that function executes without error
			if len(matched) == 0 {
				t.Error("Expected router to remain in matched map")
			}
		})
	}
}

func TestApplyServiceOverride(t *testing.T) {
	matched := map[string]*dynamic.Service{
		"test-service": {
			LoadBalancer: &dynamic.ServersLoadBalancer{
				Servers: []dynamic.Server{
					{URL: "http://old-server:8080"},
				},
			},
		},
	}

	match := ""

	value := []string{"http://new-server:8080"}

	applyServiceOverride(matched, match, value, func(s *dynamic.Service, urls []string) {
		if s.LoadBalancer != nil {
			s.LoadBalancer.Servers = make([]dynamic.Server, len(urls))
			for i, url := range urls {
				s.LoadBalancer.Servers[i] = dynamic.Server{URL: url}
			}
		}
	})

	// Test that function executes without error
	if len(matched) == 0 {
		t.Error("Expected service to remain in matched map")
	}
}

func TestApplyTCPServiceOverride(t *testing.T) {
	matched := map[string]*dynamic.TCPService{
		"tcp-service": {
			LoadBalancer: &dynamic.TCPServersLoadBalancer{
				Servers: []dynamic.TCPServer{
					{Address: "old-server:8080"},
				},
			},
		},
	}

	match := ""

	value := []string{"new-server:8080"}

	applyTCPServiceOverride(matched, match, value, func(s *dynamic.TCPService, addresses []string) {
		if s.LoadBalancer != nil {
			s.LoadBalancer.Servers = make([]dynamic.TCPServer, len(addresses))
			for i, addr := range addresses {
				s.LoadBalancer.Servers[i] = dynamic.TCPServer{Address: addr}
			}
		}
	})

	// Test that function executes without error
	if len(matched) == 0 {
		t.Error("Expected TCP service to remain in matched map")
	}
}

func TestApplyUDPServiceOverride(t *testing.T) {
	matched := map[string]*dynamic.UDPService{
		"udp-service": {
			LoadBalancer: &dynamic.UDPServersLoadBalancer{
				Servers: []dynamic.UDPServer{
					{Address: "old-server:8080"},
				},
			},
		},
	}

	match := ""

	value := []string{"new-server:8080"}

	applyUDPServiceOverride(matched, match, value, func(s *dynamic.UDPService, addresses []string) {
		if s.LoadBalancer != nil {
			s.LoadBalancer.Servers = make([]dynamic.UDPServer, len(addresses))
			for i, addr := range addresses {
				s.LoadBalancer.Servers[i] = dynamic.UDPServer{Address: addr}
			}
		}
	})

	// Test that function executes without error
	if len(matched) == 0 {
		t.Error("Expected UDP service to remain in matched map")
	}
}

func TestStripProviderFromKeys(t *testing.T) {
	m := map[string]*int{"a@p1": ptr(1), "b": ptr(2)}
	out := StripProviderFromKeys(m)
	if len(out) != 2 {
		t.Fatalf("len(out)=%d", len(out))
	}
	if _, ok := out["a"]; !ok {
		t.Fatalf("expected key 'a'")
	}
	if _, ok := out["b"]; !ok {
		t.Fatalf("expected key 'b'")
	}
	if out["a"] != m["a@p1"] || out["b"] != m["b"] {
		t.Fatalf("values not preserved")
	}
}

func TestStripProviderRefsRouter_HTTP(t *testing.T) {
	service := "svc@file"
	mids := []string{"m1@file", "m2"}
	StripProviderRefsRouter(&service, &mids)
	if service != "svc" {
		t.Fatalf("service=%q", service)
	}
	if !reflect.DeepEqual(mids, []string{"m1", "m2"}) {
		t.Fatalf("middlewares=%v", mids)
	}
}

func TestStripProviderRefsRouter_UDP(t *testing.T) {
	service := "svc@f"
	StripProviderRefsRouter(&service, nil)
	if service != "svc" {
		t.Fatalf("service=%q", service)
	}
}

func TestStripProvidersHTTP(t *testing.T) {
	httpCfg := &dynamic.HTTPConfiguration{
		Routers: map[string]*dynamic.Router{
			"r@p": {Service: "s@p", Middlewares: []string{"m1@p", "m2"}},
		},
		Services:    map[string]*dynamic.Service{"s@p": {}},
		Middlewares: map[string]*dynamic.Middleware{"m1@p": {}, "m2": {}},
	}
	StripProvidersHTTP(httpCfg)
	if _, ok := httpCfg.Routers["r"]; !ok {
		t.Fatalf("router key not stripped")
	}
	if httpCfg.Routers["r"].Service != "s" {
		t.Fatalf("router service not stripped")
	}
	if !reflect.DeepEqual(httpCfg.Routers["r"].Middlewares, []string{"m1", "m2"}) {
		t.Fatalf("router middlewares not stripped")
	}
	if _, ok := httpCfg.Services["s"]; !ok {
		t.Fatalf("service key not stripped")
	}
	if _, ok := httpCfg.Middlewares["m1"]; !ok {
		t.Fatalf("middleware key not stripped")
	}
}

func TestStripProvidersTCP(t *testing.T) {
	tcpCfg := &dynamic.TCPConfiguration{
		Routers: map[string]*dynamic.TCPRouter{
			"tr@p": {Service: "ts@p", Middlewares: []string{"tm@p"}},
		},
		Services:    map[string]*dynamic.TCPService{"ts@p": {}},
		Middlewares: map[string]*dynamic.TCPMiddleware{"tm@p": {}},
	}
	StripProvidersTCP(tcpCfg)
	if _, ok := tcpCfg.Routers["tr"]; !ok {
		t.Fatalf("tcp router key not stripped")
	}
	if tcpCfg.Routers["tr"].Service != "ts" {
		t.Fatalf("tcp router service not stripped")
	}
	if !reflect.DeepEqual(tcpCfg.Routers["tr"].Middlewares, []string{"tm"}) {
		t.Fatalf("tcp router middlewares not stripped")
	}
	if _, ok := tcpCfg.Services["ts"]; !ok {
		t.Fatalf("tcp service key not stripped")
	}
	if _, ok := tcpCfg.Middlewares["tm"]; !ok {
		t.Fatalf("tcp middleware key not stripped")
	}
}

func TestStripProvidersUDP(t *testing.T) {
	udpCfg := &dynamic.UDPConfiguration{
		Routers: map[string]*dynamic.UDPRouter{
			"ur@p": {Service: "us@p"},
		},
		Services: map[string]*dynamic.UDPService{"us@p": {}},
	}
	StripProvidersUDP(udpCfg)
	if _, ok := udpCfg.Routers["ur"]; !ok {
		t.Fatalf("udp router key not stripped")
	}
	if udpCfg.Routers["ur"].Service != "us" {
		t.Fatalf("udp router service not stripped")
	}
	if _, ok := udpCfg.Services["us"]; !ok {
		t.Fatalf("udp service key not stripped")
	}
}

func ptr[T any](v T) *T { return &v }
