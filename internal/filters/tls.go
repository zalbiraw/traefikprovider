package filters

import (
	"encoding/json"

	dynamictls "github.com/traefik/genconf/dynamic/tls"
	"github.com/zalbiraw/traefik-provider/config"
)

func TLSCertificates(certificates interface{}, config *config.TLSSection) []*dynamictls.CertAndStores {
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

func TLSOptions(options interface{}, config *config.TLSSection) map[string]dynamictls.Options {
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

func TLSStores(stores interface{}, config *config.TLSSection) map[string]dynamictls.Store {
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
