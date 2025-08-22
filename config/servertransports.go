package config

// ServerTransportsConfig matches the structure for the serverTransports section in YAML.
type ServerTransportsConfig struct {
	Discover bool `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filters  struct {
		Name      string `json:"name,omitempty" yaml:"name,omitempty"`
		NameRegex string `json:"nameRegex,omitempty" yaml:"nameRegex,omitempty"`
	} `json:"filters,omitempty" yaml:"filters,omitempty"`
	ExtraServerTransports []interface{} `json:"extraServerTransports,omitempty" yaml:"extraServerTransports,omitempty"`
}
