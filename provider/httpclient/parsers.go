package httpclient

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/traefik/genconf/dynamic"
)

// filterAndUnmarshalRouters filters routers by integer priority and unmarshals valid ones.
// filterAndUnmarshalRouters filters routers and optionally parses priority if syncPriority is true.
func filterAndUnmarshalRouters(routers interface{}, syncPriority bool) map[string]*dynamic.Router {
	result := make(map[string]*dynamic.Router)
	if routersMap, ok := routers.(map[string]interface{}); ok {
		for name, val := range routersMap {
			if routerMap, ok := val.(map[string]interface{}); ok {
				router := &dynamic.Router{}
				if syncPriority {
					if prio, ok := routerMap["priority"]; ok {
						switch v := prio.(type) {
						case float64:
							if v != float64(int(v)) || v > float64(math.MaxInt64) || v < float64(math.MinInt64) {
								fmt.Printf("[WARN] Invalid router priority for %s: %v\n", name, v)
								break
							}
							router.Priority = int(v)
						case int:
							router.Priority = v
						default:
							fmt.Printf("[WARN] Invalid router priority for %s: %v\n", name, v)
						}
					}
				}
				b, err := json.Marshal(routerMap)
				if err != nil {
					continue
				}
				if err := json.Unmarshal(b, router); err != nil {
					continue
				}
				result[name] = router
			}
		}
	}
	return result
}

func parseHTTPConfig(raw map[string]interface{}, httpConfig *dynamic.HTTPConfiguration, discoverPriority bool) error {
	if routers, ok := raw["routers"]; ok {
		httpConfig.Routers = filterAndUnmarshalRouters(routers, discoverPriority)
	}
	if services, ok := raw["services"]; ok {
		jsonServices, err := json.Marshal(services)
		if err != nil {
			return fmt.Errorf("error marshaling services: %w", err)
		}
		if err := json.Unmarshal(jsonServices, &httpConfig.Services); err != nil {
			return fmt.Errorf("error unmarshaling services: %w", err)
		}
	}
	if middlewares, ok := raw["middlewares"]; ok {
		jsonMiddlewares, err := json.Marshal(middlewares)
		if err != nil {
			return fmt.Errorf("error marshaling middlewares: %w", err)
		}
		if err := json.Unmarshal(jsonMiddlewares, &httpConfig.Middlewares); err != nil {
			return fmt.Errorf("error unmarshaling middlewares: %w", err)
		}
	}
	if serversTransports, ok := raw["serversTransports"]; ok {
		jsonST, err := json.Marshal(serversTransports)
		if err != nil {
			return fmt.Errorf("error marshaling serversTransports: %w", err)
		}
		if err := json.Unmarshal(jsonST, &httpConfig.ServersTransports); err != nil {
			return fmt.Errorf("error unmarshaling serversTransports: %w", err)
		}
	}
	return nil
}

func parseTCPConfig(raw map[string]interface{}, tcpConfig *dynamic.TCPConfiguration) error {
	if routers, ok := raw["tcpRouters"]; ok {
		jsonRouters, err := json.Marshal(routers)
		if err != nil {
			return fmt.Errorf("error marshaling tcpRouters: %w", err)
		}
		if err := json.Unmarshal(jsonRouters, &tcpConfig.Routers); err != nil {
			return fmt.Errorf("error unmarshaling tcpRouters: %w", err)
		}
	}
	if services, ok := raw["tcpServices"]; ok {
		jsonServices, err := json.Marshal(services)
		if err != nil {
			return fmt.Errorf("error marshaling tcpServices: %w", err)
		}
		if err := json.Unmarshal(jsonServices, &tcpConfig.Services); err != nil {
			return fmt.Errorf("error unmarshaling tcpServices: %w", err)
		}
	}
	if middlewares, ok := raw["tcpMiddlewares"]; ok {
		jsonMiddlewares, err := json.Marshal(middlewares)
		if err != nil {
			return fmt.Errorf("error marshaling tcpMiddlewares: %w", err)
		}
		if err := json.Unmarshal(jsonMiddlewares, &tcpConfig.Middlewares); err != nil {
			return fmt.Errorf("error unmarshaling tcpMiddlewares: %w", err)
		}
	}
	return nil
}

func parseUDPConfig(raw map[string]interface{}, udpConfig *dynamic.UDPConfiguration) error {
	if routers, ok := raw["udpRouters"]; ok {
		jsonRouters, err := json.Marshal(routers)
		if err != nil {
			return fmt.Errorf("error marshaling udpRouters: %w", err)
		}
		if err := json.Unmarshal(jsonRouters, &udpConfig.Routers); err != nil {
			return fmt.Errorf("error unmarshaling udpRouters: %w", err)
		}
	}
	if services, ok := raw["udpServices"]; ok {
		jsonServices, err := json.Marshal(services)
		if err != nil {
			return fmt.Errorf("error marshaling udpServices: %w", err)
		}
		if err := json.Unmarshal(jsonServices, &udpConfig.Services); err != nil {
			return fmt.Errorf("error unmarshaling udpServices: %w", err)
		}
	}
	return nil
}

func parseTLSConfig(raw map[string]interface{}, tlsConfig *dynamic.TLSConfiguration) error {
	if certificates, ok := raw["tlsCertificates"]; ok {
		jsonCerts, err := json.Marshal(certificates)
		if err != nil {
			return fmt.Errorf("error marshaling tlsCertificates: %w", err)
		}
		if err := json.Unmarshal(jsonCerts, &tlsConfig.Certificates); err != nil {
			return fmt.Errorf("error unmarshaling tlsCertificates: %w", err)
		}
	}
	if options, ok := raw["tlsOptions"]; ok {
		jsonOptions, err := json.Marshal(options)
		if err != nil {
			return fmt.Errorf("error marshaling tlsOptions: %w", err)
		}
		if err := json.Unmarshal(jsonOptions, &tlsConfig.Options); err != nil {
			return fmt.Errorf("error unmarshaling tlsOptions: %w", err)
		}
	}
	if stores, ok := raw["tlsStores"]; ok {
		jsonStores, err := json.Marshal(stores)
		if err != nil {
			return fmt.Errorf("error marshaling tlsStores: %w", err)
		}
		if err := json.Unmarshal(jsonStores, &tlsConfig.Stores); err != nil {
			return fmt.Errorf("error unmarshaling tlsStores: %w", err)
		}
	}
	return nil
}
