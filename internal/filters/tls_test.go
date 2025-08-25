package filters

import (
	"testing"
)

func TestTLSCertificates(t *testing.T) {
	tests := []struct {
		name         string
		certificates interface{}
		expectedLen  int
	}{
		{
			name: "filter all certificates",
			certificates: []interface{}{
				map[string]interface{}{"certFile": "/path/to/cert1.pem", "keyFile": "/path/to/key1.pem"},
				map[string]interface{}{"certFile": "/path/to/cert2.pem", "keyFile": "/path/to/key2.pem"},
			},
			expectedLen: 2,
		},
		{
			name:         "empty certificates",
			certificates: []interface{}{},
			expectedLen:  0,
		},
		{
			name:         "invalid certificates type",
			certificates: "invalid",
			expectedLen:  0,
		},
		{
			name: "invalid certificate in slice",
			certificates: []interface{}{
				map[string]interface{}{"certFile": "/path/to/cert1.pem", "keyFile": "/path/to/key1.pem"},
				"invalid",
			},
			expectedLen: 1,
		},
		{
			name: "unmarshalable certificate",
			certificates: []interface{}{
				map[string]interface{}{"invalid": make(chan int)},
			},
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TLSCertificates(tt.certificates, nil)

			if len(result) != tt.expectedLen {
				t.Errorf("Expected %d certificates, got %d", tt.expectedLen, len(result))
			}
		})
	}
}

func TestTLSOptions(t *testing.T) {
	tests := []struct {
		name        string
		options     interface{}
		expectedLen int
	}{
		{
			name: "filter all options",
			options: map[string]interface{}{
				"default": map[string]interface{}{"minVersion": "VersionTLS12"},
				"strict":  map[string]interface{}{"minVersion": "VersionTLS13"},
			},
			expectedLen: 2,
		},
		{
			name:        "empty options",
			options:     map[string]interface{}{},
			expectedLen: 0,
		},
		{
			name:        "invalid options type",
			options:     "invalid",
			expectedLen: 0,
		},
		{
			name: "invalid option in map",
			options: map[string]interface{}{
				"valid":   map[string]interface{}{"minVersion": "VersionTLS12"},
				"invalid": "not-a-map",
			},
			expectedLen: 1,
		},
		{
			name: "unmarshalable option",
			options: map[string]interface{}{
				"invalid": map[string]interface{}{"channel": make(chan int)},
			},
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TLSOptions(tt.options, nil)

			if len(result) != tt.expectedLen {
				t.Errorf("Expected %d options, got %d", tt.expectedLen, len(result))
			}
		})
	}
}

func TestTLSStores(t *testing.T) {
	tests := []struct {
		name        string
		stores      interface{}
		expectedLen int
	}{
		{
			name: "filter all stores",
			stores: map[string]interface{}{
				"default": map[string]interface{}{"defaultCertificate": map[string]interface{}{"certFile": "cert.pem", "keyFile": "key.pem"}},
				"custom":  map[string]interface{}{"defaultCertificate": map[string]interface{}{"certFile": "custom.pem", "keyFile": "custom-key.pem"}},
			},
			expectedLen: 2,
		},
		{
			name:        "empty stores",
			stores:      map[string]interface{}{},
			expectedLen: 0,
		},
		{
			name:        "invalid stores type",
			stores:      "invalid",
			expectedLen: 0,
		},
		{
			name: "invalid store in map",
			stores: map[string]interface{}{
				"valid":   map[string]interface{}{"defaultCertificate": map[string]interface{}{"certFile": "cert.pem"}},
				"invalid": "not-a-map",
			},
			expectedLen: 1,
		},
		{
			name: "unmarshalable store",
			stores: map[string]interface{}{
				"invalid": map[string]interface{}{"channel": make(chan int)},
			},
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TLSStores(tt.stores, nil)

			if len(result) != tt.expectedLen {
				t.Errorf("Expected %d stores, got %d", tt.expectedLen, len(result))
			}
		})
	}
}

func TestTLSInvalidData(t *testing.T) {
	// Test with invalid data that causes JSON marshal errors
	invalidData := map[string]interface{}{
		"test": map[string]interface{}{
			"invalid": make(chan int), // Cannot be marshaled to JSON
		},
	}

	// These should handle marshal errors gracefully
	certsResult := TLSCertificates(invalidData, nil)
	if len(certsResult) != 0 {
		t.Error("Expected empty result for unmarshalable TLS certificates")
	}

	optionsResult := TLSOptions(invalidData, nil)
	if len(optionsResult) != 0 {
		t.Error("Expected empty result for unmarshalable TLS options")
	}

	storesResult := TLSStores(invalidData, nil)
	if len(storesResult) != 0 {
		t.Error("Expected empty result for unmarshalable TLS stores")
	}
}
