package config

// ServicesConfig holds discovery and override settings for services.
type ServicesConfig struct {
	Discover      bool             `json:"discover,omitempty" yaml:"discover,omitempty"`
	Matcher       string           `json:"matcher,omitempty" yaml:"matcher,omitempty"`
	Overrides     ServiceOverrides `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraServices []interface{}    `json:"extraServices,omitempty" yaml:"extraServices,omitempty"`
}

// ServiceOverrides defines how to override service backends and healthchecks.
type ServiceOverrides struct {
	Servers      []OverrideServer      `json:"servers,omitempty" yaml:"servers,omitempty"`
	Healthchecks []OverrideHealthcheck `json:"healthchecks,omitempty" yaml:"healthchecks,omitempty"`
}

// OverrideServer configures server address overrides for matching services.
type OverrideServer struct {
	Strategy string      `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Value    interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	Matcher  string      `json:"matcher,omitempty" yaml:"matcher,omitempty"`
}

// OverrideHealthcheck overrides healthcheck settings for matching services.
type OverrideHealthcheck struct {
	Path     string `json:"path,omitempty" yaml:"path,omitempty"`
	Interval string `json:"interval,omitempty" yaml:"interval,omitempty"`
	Timeout  string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Matcher  string `json:"matcher,omitempty" yaml:"matcher,omitempty"`
}
