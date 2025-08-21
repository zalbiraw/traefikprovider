package config

type MiddlewaresConfig struct {
	Discover        bool                 `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filters         MiddlewareFilters    `json:"filters,omitempty" yaml:"filters,omitempty"`
	ExtraMiddlewares []interface{}       `json:"extraMiddlewares,omitempty" yaml:"extraMiddlewares,omitempty"`
}

type MiddlewareFilters struct {
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	NameRegex string `json:"nameRegex,omitempty" yaml:"nameRegex,omitempty"`
}

type ServicesConfig struct {
	Discover      bool                `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filters       ServiceFilters      `json:"filters,omitempty" yaml:"filters,omitempty"`
	Overrides     ServiceOverrides   `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraServices []interface{}      `json:"extraServices,omitempty" yaml:"extraServices,omitempty"`
}

type ServiceFilters struct {
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	NameRegex string `json:"nameRegex,omitempty" yaml:"nameRegex,omitempty"`
}

type ServiceOverrides struct {
	Servers      []OverrideServer      `json:"servers,omitempty" yaml:"servers,omitempty"`
	Healthchecks []OverrideHealthcheck `json:"healthchecks,omitempty" yaml:"healthchecks,omitempty"`
}

type OverrideServer struct {
	Mode      string                 `json:"mode,omitempty" yaml:"mode,omitempty"`
	Strategy  string                 `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Values    []string               `json:"values,omitempty" yaml:"values,omitempty"`
	Filters   map[string]interface{}  `json:"filters,omitempty" yaml:"filters,omitempty"`
	Connection *ConnectionConfig      `json:"connection,omitempty" yaml:"connection,omitempty"`
}

type OverrideHealthcheck struct {
	Mode      string                 `json:"mode,omitempty" yaml:"mode,omitempty"`
	Path      string                 `json:"path,omitempty" yaml:"path,omitempty"`
	PathRegex string                 `json:"pathRegex,omitempty" yaml:"pathRegex,omitempty"`
	Interval  string                 `json:"interval,omitempty" yaml:"interval,omitempty"`
	Timeout   string                 `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Filters   map[string]interface{}  `json:"filters,omitempty" yaml:"filters,omitempty"`
}
