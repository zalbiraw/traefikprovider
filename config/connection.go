package config

type ConnectionConfig struct {
	Host    []string          `json:"host,omitempty" yaml:"host,omitempty"`
	Port    int               `json:"port,omitempty" yaml:"port,omitempty"`
	Path    string            `json:"path,omitempty" yaml:"path,omitempty"`
	Timeout string            `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	MTLS    *MTLSConfig       `json:"mTLS,omitempty" yaml:"mTLS,omitempty"`
}

type MTLSConfig struct {
	CAFile   string `json:"caFile,omitempty" yaml:"caFile,omitempty"`
	CertFile string `json:"certFile,omitempty" yaml:"certFile,omitempty"`
	KeyFile  string `json:"keyFile,omitempty" yaml:"keyFile,omitempty"`
}
