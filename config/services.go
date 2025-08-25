package config

// ServicesConfig holds discovery and override settings for services.
type ServicesConfig struct {
	Discover      bool             `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filter        ServiceFilter    `json:"filter,omitempty" yaml:"filter,omitempty"`
	Overrides     ServiceOverrides `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraServices []interface{}    `json:"extraServices,omitempty" yaml:"extraServices,omitempty"`
}

// ServiceFilter filters services by name and provider.
type ServiceFilter struct {
	Name     string `json:"name,omitempty" yaml:"name,omitempty"`
	Provider string `json:"provider,omitempty" yaml:"provider,omitempty"`
}

// ServiceOverrides defines how to override service backends and healthchecks.
type ServiceOverrides struct {
	Servers      []OverrideServer      `json:"servers,omitempty" yaml:"servers,omitempty"`
	Healthchecks []OverrideHealthcheck `json:"healthchecks,omitempty" yaml:"healthchecks,omitempty"`
}

// OverrideServer configures server address overrides for matching services.
type OverrideServer struct {
	Strategy string        `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Value    interface{}   `json:"value,omitempty" yaml:"value,omitempty"`
	Filter   ServiceFilter `json:"filter,omitempty" yaml:"filter,omitempty"`
	Tunnel   string        `json:"tunnel,omitempty" yaml:"tunnel,omitempty"`
}

// OverrideHealthcheck overrides healthcheck settings for matching services.
type OverrideHealthcheck struct {
	Path     string        `json:"path,omitempty" yaml:"path,omitempty"`
	Interval string        `json:"interval,omitempty" yaml:"interval,omitempty"`
	Timeout  string        `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Filter   ServiceFilter `json:"filter,omitempty" yaml:"filter,omitempty"`
}
