package overrides

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
	"github.com/zalbiraw/traefikprovider/internal/filters"
)

func applyRouterOverride[T any](filtered map[string]*dynamic.Router, routerFilter config.RouterFilter, value T, apply func(r *dynamic.Router, v T)) {
	rc := &config.RoutersConfig{Filter: routerFilter}
	for key, router := range filters.HTTPRouters(filtered, rc, config.ProviderFilter{}) {
		apply(router, value)
		filtered[key] = router
	}
}

func handleRouterOverride(
	filtered map[string]*dynamic.Router,
	routerFilter config.RouterFilter,
	value interface{},
	applyArray func(r *dynamic.Router, arr []string),
	applyString func(r *dynamic.Router, s string),
) {
	switch v := value.(type) {
	case []string:
		applyRouterOverride(filtered, routerFilter, v, applyArray)
	case string:
		applyRouterOverride(filtered, routerFilter, v, applyString)
	}
}

func applyServiceOverride[T any](filtered map[string]*dynamic.Service, serviceFilter config.ServiceFilter, value T, apply func(r *dynamic.Service, v T)) {
	rc := &config.ServicesConfig{Filter: serviceFilter}
	for key, service := range filters.HTTPServices(filtered, rc, config.ProviderFilter{}) {
		apply(service, value)
		filtered[key] = service
	}
}

func handleServiceOverride(
	filtered map[string]*dynamic.Service,
	serviceFilter config.ServiceFilter,
	value interface{},
	applyArray func(r *dynamic.Service, arr []string),
	applyString func(r *dynamic.Service, s string),
) {
	switch v := value.(type) {
	case []string:
		applyServiceOverride(filtered, serviceFilter, v, applyArray)
	case string:
		applyServiceOverride(filtered, serviceFilter, v, applyString)
	}
}

func applyTCPServiceOverride[T any](filtered map[string]*dynamic.TCPService, serviceFilter config.ServiceFilter, value T, apply func(r *dynamic.TCPService, v T)) {
	rc := &config.ServicesConfig{Filter: serviceFilter}
	for key, service := range filters.TCPServices(filtered, rc, config.ProviderFilter{}) {
		apply(service, value)
		filtered[key] = service
	}
}

func handleTCPServiceOverride(
	filtered map[string]*dynamic.TCPService,
	serviceFilter config.ServiceFilter,
	value interface{},
	applyArray func(r *dynamic.TCPService, arr []string),
	applyString func(r *dynamic.TCPService, s string),
) {
	switch v := value.(type) {
	case []string:
		applyTCPServiceOverride(filtered, serviceFilter, v, applyArray)
	case string:
		applyTCPServiceOverride(filtered, serviceFilter, v, applyString)
	}
}

func applyUDPServiceOverride[T any](filtered map[string]*dynamic.UDPService, serviceFilter config.ServiceFilter, value T, apply func(r *dynamic.UDPService, v T)) {
	rc := &config.UDPServicesConfig{Filter: serviceFilter}
	for key, service := range filters.UDPServices(filtered, rc, config.ProviderFilter{}) {
		apply(service, value)
		filtered[key] = service
	}
}

func handleUDPServiceOverride(
	filtered map[string]*dynamic.UDPService,
	serviceFilter config.ServiceFilter,
	value interface{},
	applyArray func(r *dynamic.UDPService, arr []string),
	applyString func(r *dynamic.UDPService, s string),
) {
	switch v := value.(type) {
	case []string:
		applyUDPServiceOverride(filtered, serviceFilter, v, applyArray)
	case string:
		applyUDPServiceOverride(filtered, serviceFilter, v, applyString)
	}
}
