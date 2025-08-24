package config

type ServicesConfig struct {
	Discover      bool             `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filters       ServiceFilters   `json:"filters,omitempty" yaml:"filters,omitempty"`
	Overrides     ServiceOverrides `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraServices []interface{}    `json:"extraServices,omitempty" yaml:"extraServices,omitempty"`
}

type ServiceFilters struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type ServiceOverrides struct {
	Servers      []OverrideServer      `json:"servers,omitempty" yaml:"servers,omitempty"`
	Healthchecks []OverrideHealthcheck `json:"healthchecks,omitempty" yaml:"healthchecks,omitempty"`
}

type OverrideServer struct {
	Strategy string                 `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Value    []string               `json:"value,omitempty" yaml:"value,omitempty"`
	Filters  map[string]interface{} `json:"filters,omitempty" yaml:"filters,omitempty"`
	Tunnel   string                 `json:"tunnel,omitempty" yaml:"tunnel,omitempty"`
}

type OverrideHealthcheck struct {
	Path     string                 `json:"path,omitempty" yaml:"path,omitempty"`
	Interval string                 `json:"interval,omitempty" yaml:"interval,omitempty"`
	Timeout  string                 `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Filters  map[string]interface{} `json:"filters,omitempty" yaml:"filters,omitempty"`
}
