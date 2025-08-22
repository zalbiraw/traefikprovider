package httpclient

import (
	"encoding/json"

dynamictls "github.com/traefik/genconf/dynamic/tls"
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

// filterCertificates filters TLS certificates based on the provided config.
func filterCertificates(certificates interface{}, config interface{}) []dynamictls.Certificate {
	certs := []dynamictls.Certificate{}
	certsSlice, ok := certificates.([]interface{})
	if !ok {
		return certs
	}
	for _, val := range certsSlice {
		certMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		cert := dynamictls.Certificate{}
		b, err := json.Marshal(certMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, &cert); err != nil {
			continue
		}
		certs = append(certs, cert)
	}
	return certs
}

// filterOptions filters TLS options based on the provided config.
func filterOptions(options interface{}, config interface{}) map[string]*dynamictls.Options {
	result := make(map[string]*dynamictls.Options)
	optionsMap, ok := options.(map[string]interface{})
	if !ok {
		return result
	}
	for name, val := range optionsMap {
		optMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		opt := &dynamictls.Options{}
		b, err := json.Marshal(optMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, opt); err != nil {
			continue
		}
		result[name] = opt
	}
	return result
}

// filterTLSCertificates filters TLS certificates for TLSConfiguration.
func filterTLSCertificates(certificates interface{}, config *config.TLSSection) []*dynamictls.CertAndStores {
	certs := []*dynamictls.CertAndStores{}
	certsSlice, ok := certificates.([]interface{})
	if !ok {
		return certs
	}
	for _, val := range certsSlice {
		certMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		cert := dynamictls.Certificate{}
		b, err := json.Marshal(certMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, &cert); err != nil {
			continue
		}
		certs = append(certs, &dynamictls.CertAndStores{Certificate: cert})
	}
	return certs
}

// filterTLSOptions filters TLS options for TLSConfiguration.
func filterTLSOptions(options interface{}, config *config.TLSSection) map[string]dynamictls.Options {
	result := make(map[string]dynamictls.Options)
	optionsMap, ok := options.(map[string]interface{})
	if !ok {
		return result
	}
	for name, val := range optionsMap {
		optMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		opt := dynamictls.Options{}
		b, err := json.Marshal(optMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, &opt); err != nil {
			continue
		}
		result[name] = opt
	}
	return result
}

// filterTLSStores filters TLS stores for TLSConfiguration.
func filterTLSStores(stores interface{}, config *config.TLSSection) map[string]dynamictls.Store {
	result := make(map[string]dynamictls.Store)
	storesMap, ok := stores.(map[string]interface{})
	if !ok {
		return result
	}
	for name, val := range storesMap {
		storeMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		store := dynamictls.Store{}
		b, err := json.Marshal(storeMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, &store); err != nil {
			continue
		}
		result[name] = store
	}
	return result
}

// filterStores filters TLS stores based on the provided config.
func filterStores(stores interface{}, config interface{}) map[string]*dynamictls.Store {
	result := make(map[string]*dynamictls.Store)
	storesMap, ok := stores.(map[string]interface{})
	if !ok {
		return result
	}
	for name, val := range storesMap {
		storeMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		store := &dynamictls.Store{}
		b, err := json.Marshal(storeMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, store); err != nil {
			continue
		}
		result[name] = store
	}
	return result
}

// filterTCPRouters filters TCP routers based on the provided config.
func filterTCPRouters(routers interface{}, config *config.RoutersConfig) map[string]*dynamic.TCPRouter {
	result := make(map[string]*dynamic.TCPRouter)
	routersMap, ok := routers.(map[string]interface{})
	if !ok {
		return result
	}
	for name, val := range routersMap {
		routerMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		router := &dynamic.TCPRouter{}
		b, err := json.Marshal(routerMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, router); err != nil {
			continue
		}
		result[name] = router
	}
	return result
}

// filterTCPServices filters TCP services based on the provided config.
func filterTCPServices(services interface{}, config *config.ServicesConfig) map[string]*dynamic.TCPService {
	result := make(map[string]*dynamic.TCPService)
	servicesMap, ok := services.(map[string]interface{})
	if !ok {
		return result
	}
	for name, val := range servicesMap {
		serviceMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		service := &dynamic.TCPService{}
		b, err := json.Marshal(serviceMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, service); err != nil {
			continue
		}
		result[name] = service
	}
	return result
}

// filterTCPMiddlewares filters TCP middlewares based on the provided config.
func filterTCPMiddlewares(middlewares interface{}, config *config.MiddlewaresConfig) map[string]*dynamic.TCPMiddleware {
	result := make(map[string]*dynamic.TCPMiddleware)
	middlewaresMap, ok := middlewares.(map[string]interface{})
	if !ok {
		return result
	}
	for name, val := range middlewaresMap {
		middlewareMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		middleware := &dynamic.TCPMiddleware{}
		b, err := json.Marshal(middlewareMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, middleware); err != nil {
			continue
		}
		result[name] = middleware
	}
	return result
}

// filterUDPRouters filters UDP routers based on the provided config.
func filterUDPRouters(routers interface{}, config *config.UDPRoutersConfig) map[string]*dynamic.UDPRouter {
	result := make(map[string]*dynamic.UDPRouter)
	routersMap, ok := routers.(map[string]interface{})
	if !ok {
		return result
	}
	for name, val := range routersMap {
		routerMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		router := &dynamic.UDPRouter{}
		b, err := json.Marshal(routerMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, router); err != nil {
			continue
		}
		result[name] = router
	}
	return result
}

// filterUDPServices filters UDP services based on the provided config.
func filterUDPServices(services interface{}, config *config.UDPServicesConfig) map[string]*dynamic.UDPService {
	result := make(map[string]*dynamic.UDPService)
	servicesMap, ok := services.(map[string]interface{})
	if !ok {
		return result
	}
	for name, val := range servicesMap {
		serviceMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		service := &dynamic.UDPService{}
		b, err := json.Marshal(serviceMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, service); err != nil {
			continue
		}
		result[name] = service
	}
	return result
}

// filterRouters filters routers based on the provided config.
func filterRouters(routers interface{}, config *config.RoutersConfig) map[string]*dynamic.Router {
	result := make(map[string]*dynamic.Router)
	routersMap, ok := routers.(map[string]interface{})
	if !ok {
		return result
	}
	for name, val := range routersMap {
		routerMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		router := &dynamic.Router{}
		if config.DiscoverPriority {
			router.Priority = extractRouterPriority(routerMap, name)
		}
		if err := unmarshalRouter(routerMap, router); err != nil {
			continue
		}
		result[name] = router
	}
	return result
}

// filterServices filters services based on the provided config.
func filterServices(services interface{}, config *config.ServicesConfig) map[string]*dynamic.Service {
	result := make(map[string]*dynamic.Service)
	servicesMap, ok := services.(map[string]interface{})
	if !ok {
		return result
	}
	for name, val := range servicesMap {
		serviceMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		service := &dynamic.Service{}
		// Add filtering logic here based on config.Filters if needed
		// For now, just add all
		// TODO: implement actual filtering based on config
		b, err := json.Marshal(serviceMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, service); err != nil {
			continue
		}
		result[name] = service
	}
	return result
}

// filterMiddlewares filters middlewares based on the provided config.
func filterMiddlewares(middlewares interface{}, config *config.MiddlewaresConfig) map[string]*dynamic.Middleware {
	result := make(map[string]*dynamic.Middleware)
	middlewaresMap, ok := middlewares.(map[string]interface{})
	if !ok {
		return result
	}
	for name, val := range middlewaresMap {
		middlewareMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		middleware := &dynamic.Middleware{}
		// Add filtering logic here based on config.Filters if needed
		// TODO: implement actual filtering based on config
		b, err := json.Marshal(middlewareMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, middleware); err != nil {
			continue
		}
		result[name] = middleware
	}
	return result
}

// filterServerTransports filters server transports based on the provided config.
func filterServerTransports(serverTransports interface{}, config *config.ServerTransportsConfig) map[string]*dynamic.ServersTransport {
	result := make(map[string]*dynamic.ServersTransport)
	stMap, ok := serverTransports.(map[string]interface{})
	if !ok {
		return result
	}
	for name, val := range stMap {
		stItemMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		st := &dynamic.ServersTransport{}
		// Add filtering logic here based on config.Filters if needed
		// TODO: implement actual filtering based on config
		b, err := json.Marshal(stItemMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, st); err != nil {
			continue
		}
		result[name] = st
	}
	return result
}
