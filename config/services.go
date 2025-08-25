package config

type ServicesConfig struct {
	Discover      bool             `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filter        ServiceFilter    `json:"filter,omitempty" yaml:"filter,omitempty"`
	Overrides     ServiceOverrides `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraServices []interface{}    `json:"extraServices,omitempty" yaml:"extraServices,omitempty"`
}

type ServiceFilter struct {
	Name     string `json:"name,omitempty" yaml:"name,omitempty"`
	Provider string `json:"provider,omitempty" yaml:"provider,omitempty"`
}

type ServiceOverrides struct {
	Servers      []OverrideServer      `json:"servers,omitempty" yaml:"servers,omitempty"`
	Healthchecks []OverrideHealthcheck `json:"healthchecks,omitempty" yaml:"healthchecks,omitempty"`
}

type OverrideServer struct {
	Strategy string        `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Value    interface{}   `json:"value,omitempty" yaml:"value,omitempty"`
	Filter   ServiceFilter `json:"filter,omitempty" yaml:"filter,omitempty"`
	Tunnel   string        `json:"tunnel,omitempty" yaml:"tunnel,omitempty"`
}

type OverrideHealthcheck struct {
	Path     string        `json:"path,omitempty" yaml:"path,omitempty"`
	Interval string        `json:"interval,omitempty" yaml:"interval,omitempty"`
	Timeout  string        `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Filter   ServiceFilter `json:"filter,omitempty" yaml:"filter,omitempty"`
}
