package overrides

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
	"github.com/zalbiraw/traefik-provider/internal/filters"
)

func applyRouterOverride[T any](filtered map[string]*dynamic.Router, routerFilters config.RouterFilters, value T, apply func(r *dynamic.Router, v T)) {
	rc := &config.RoutersConfig{Filters: routerFilters}
	for key, router := range filters.HTTPRouters(filtered, rc) {
		apply(router, value)
		filtered[key] = router
	}
}

func handleRouterOverride(
	filtered map[string]*dynamic.Router,
	routerFilters config.RouterFilters,
	value interface{},
	applyArray func(r *dynamic.Router, arr []string),
	applyString func(r *dynamic.Router, s string),
) {
	switch v := value.(type) {
	case []string:
		applyRouterOverride(filtered, routerFilters, v, applyArray)
	case string:
		applyRouterOverride(filtered, routerFilters, v, applyString)
	}
}

func applyServiceOverride[T any](filtered map[string]*dynamic.Service, serviceFilters config.ServiceFilters, value T, apply func(r *dynamic.Service, v T)) {
	rc := &config.ServicesConfig{Filters: serviceFilters}
	for key, service := range filters.HTTPServices(filtered, rc) {
		apply(service, value)
		filtered[key] = service
	}
}

func handleServiceOverride(
	filtered map[string]*dynamic.Service,
	serviceFilters config.ServiceFilters,
	value interface{},
	applyArray func(r *dynamic.Service, arr []string),
	applyString func(r *dynamic.Service, s string),
) {
	switch v := value.(type) {
	case []string:
		applyServiceOverride(filtered, serviceFilters, v, applyArray)
	case string:
		applyServiceOverride(filtered, serviceFilters, v, applyString)
	}
}

func applyTCPServiceOverride[T any](filtered map[string]*dynamic.TCPService, serviceFilters config.ServiceFilters, value T, apply func(r *dynamic.TCPService, v T)) {
	rc := &config.ServicesConfig{Filters: serviceFilters}
	for key, service := range filters.TCPServices(filtered, rc) {
		apply(service, value)
		filtered[key] = service
	}
}

func handleTCPServiceOverride(
	filtered map[string]*dynamic.TCPService,
	serviceFilters config.ServiceFilters,
	value interface{},
	applyArray func(r *dynamic.TCPService, arr []string),
	applyString func(r *dynamic.TCPService, s string),
) {
	switch v := value.(type) {
	case []string:
		applyTCPServiceOverride(filtered, serviceFilters, v, applyArray)
	case string:
		applyTCPServiceOverride(filtered, serviceFilters, v, applyString)
	}
}

func applyUDPServiceOverride[T any](filtered map[string]*dynamic.UDPService, serviceFilters config.ServiceFilters, value T, apply func(r *dynamic.UDPService, v T)) {
	rc := &config.UDPServicesConfig{Filters: serviceFilters}
	for key, service := range filters.UDPServices(filtered, rc) {
		apply(service, value)
		filtered[key] = service
	}
}

func handleUDPServiceOverride(
	filtered map[string]*dynamic.UDPService,
	serviceFilters config.ServiceFilters,
	value interface{},
	applyArray func(r *dynamic.UDPService, arr []string),
	applyString func(r *dynamic.UDPService, s string),
) {
	switch v := value.(type) {
	case []string:
		applyUDPServiceOverride(filtered, serviceFilters, v, applyArray)
	case string:
		applyUDPServiceOverride(filtered, serviceFilters, v, applyString)
	}
}
