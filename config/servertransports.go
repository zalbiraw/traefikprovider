package config

// ServerTransportsConfig matches the structure for the serverTransports section in YAML.
type ServerTransportsConfig struct {
	Discover              bool                   `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filters               ServerTransportFilters `json:"filters,omitempty" yaml:"filters,omitempty"`
	ExtraServerTransports []interface{}          `json:"extraServerTransports,omitempty" yaml:"extraServerTransports,omitempty"`
}

type ServerTransportFilters struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}
