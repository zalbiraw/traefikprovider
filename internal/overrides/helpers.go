package overrides

import (
	"encoding/json"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
	"github.com/zalbiraw/traefik-provider/internal/filters"
)

func applyRouterOverride[T any](filtered map[string]*dynamic.Router, filtersMap map[string]interface{}, value T, apply func(r *dynamic.Router, v T)) {
	var rf config.RouterFilters
	b, _ := json.Marshal(filtersMap)
	_ = json.Unmarshal(b, &rf)
	rc := &config.RoutersConfig{Filters: rf}
	for key, router := range filters.HTTPRouters(filtered, rc) {
		apply(router, value)
		filtered[key] = router
	}
}

func handleRouterOverride(
	filtered map[string]*dynamic.Router,
	filtersMap map[string]interface{},
	value interface{},
	applyArray func(r *dynamic.Router, arr []string),
	applyString func(r *dynamic.Router, s string),
) {
	switch v := value.(type) {
	case []string:
		applyRouterOverride(filtered, filtersMap, v, applyArray)
	case string:
		applyRouterOverride(filtered, filtersMap, v, applyString)
	}
}

func applyServiceOverride[T any](filtered map[string]*dynamic.Service, filtersMap map[string]interface{}, value T, apply func(r *dynamic.Service, v T)) {
	var sf config.ServiceFilters
	b, _ := json.Marshal(filtersMap)
	_ = json.Unmarshal(b, &sf)
	rc := &config.ServicesConfig{Filters: sf}
	for key, service := range filters.HTTPServices(filtered, rc) {
		apply(service, value)
		filtered[key] = service
	}
}

func handleServiceOverride(
	filtered map[string]*dynamic.Service,
	filtersMap map[string]interface{},
	value interface{},
	applyArray func(r *dynamic.Service, arr []string),
	applyString func(r *dynamic.Service, s string),
) {
	switch v := value.(type) {
	case []string:
		applyServiceOverride(filtered, filtersMap, v, applyArray)
	case string:
		applyServiceOverride(filtered, filtersMap, v, applyString)
	}
}

func applyTCPServiceOverride[T any](filtered map[string]*dynamic.TCPService, filtersMap map[string]interface{}, value T, apply func(r *dynamic.TCPService, v T)) {
	var sf config.ServiceFilters
	b, _ := json.Marshal(filtersMap)
	_ = json.Unmarshal(b, &sf)
	rc := &config.ServicesConfig{Filters: sf}
	for key, service := range filters.TCPServices(filtered, rc) {
		apply(service, value)
		filtered[key] = service
	}
}

func handleTCPServiceOverride(
	filtered map[string]*dynamic.TCPService,
	filtersMap map[string]interface{},
	value interface{},
	applyArray func(r *dynamic.TCPService, arr []string),
	applyString func(r *dynamic.TCPService, s string),
) {
	switch v := value.(type) {
	case []string:
		applyTCPServiceOverride(filtered, filtersMap, v, applyArray)
	case string:
		applyTCPServiceOverride(filtered, filtersMap, v, applyString)
	}
}

func applyUDPServiceOverride[T any](filtered map[string]*dynamic.UDPService, filtersMap map[string]interface{}, value T, apply func(r *dynamic.UDPService, v T)) {
	var sf config.ServiceFilters
	b, _ := json.Marshal(filtersMap)
	_ = json.Unmarshal(b, &sf)
	rc := &config.UDPServicesConfig{Filters: sf}
	for key, service := range filters.UDPServices(filtered, rc) {
		apply(service, value)
		filtered[key] = service
	}
}

func handleUDPServiceOverride(
	filtered map[string]*dynamic.UDPService,
	filtersMap map[string]interface{},
	value interface{},
	applyArray func(r *dynamic.UDPService, arr []string),
	applyString func(r *dynamic.UDPService, s string),
) {
	switch v := value.(type) {
	case []string:
		applyUDPServiceOverride(filtered, filtersMap, v, applyArray)
	case string:
		applyUDPServiceOverride(filtered, filtersMap, v, applyString)
	}
}
