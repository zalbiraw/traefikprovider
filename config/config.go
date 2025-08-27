// Package config defines the provider configuration model and parsing helpers.
package config

// ProviderConfig defines a single upstream provider to poll and filter.
type ProviderConfig struct {
	Name       string           `json:"name,omitempty" yaml:"name,omitempty"`
	Matcher    string           `json:"matcher,omitempty" yaml:"matcher,omitempty"`
	Connection ConnectionConfig `json:"connection,omitempty" yaml:"connection,omitempty"`
	HTTP       *HTTPSection     `json:"http,omitempty" yaml:"http,omitempty"`
	TCP        *TCPSection      `json:"tcp,omitempty" yaml:"tcp,omitempty"`
	UDP        *UDPSection      `json:"udp,omitempty" yaml:"udp,omitempty"`
	TLS        *TLSSection      `json:"tls,omitempty" yaml:"tls,omitempty"`
	Tunnels    []TunnelConfig   `json:"tunnels,omitempty" yaml:"tunnels,omitempty"`
}

// ConnectionConfig configures how to connect to the upstream provider API.
type ConnectionConfig struct {
	Host    string            `json:"host,omitempty" yaml:"host,omitempty"`
	Port    int               `json:"port,omitempty" yaml:"port,omitempty"`
	Path    string            `json:"path,omitempty" yaml:"path,omitempty"`
	Timeout string            `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	MTLS    *MTLSConfig       `json:"mTLS,omitempty" yaml:"mTLS,omitempty"` //nolint:tagliatelle
}

// MTLSConfig holds mutual TLS file paths for establishing mTLS connections.
type MTLSConfig struct {
	CAFile   string `json:"caFile,omitempty" yaml:"caFile,omitempty"`
	CertFile string `json:"certFile,omitempty" yaml:"certFile,omitempty"`
	KeyFile  string `json:"keyFile,omitempty" yaml:"keyFile,omitempty"`
	// StripRouterTLSOptions controls whether to strip router TLS options when
	// a router uses a service matched by an mTLS-enabled tunnel. When true (default),
	// if a router's TLS config only has Options set, the TLS block is removed.
	// Otherwise, only the Options field is cleared. Set to false to disable.
	StripRouterTLSOptions *bool `json:"stripRouterTLSOptions,omitempty" yaml:"stripRouterTLSOptions,omitempty"` //nolint:tagliatelle
}

// TunnelConfig defines a list of addresses grouped under a name, optionally with mTLS.
type TunnelConfig struct {
	Addresses []string    `json:"connection,omitempty" yaml:"connection,omitempty"`
	MTLS      *MTLSConfig `json:"mTLS,omitempty" yaml:"mTLS,omitempty"` //nolint:tagliatelle
	Matcher   string      `json:"matcher,omitempty" yaml:"matcher,omitempty"`
}
