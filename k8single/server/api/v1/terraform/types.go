package terraform

type State struct {
	Version          uint64            `json:"version"`
	TerraformVersion string            `json:"terraform_version"`
	Serial           uint64            `json:"serial"`
	Lineage          string            `json:"lineage"`
	Outputs          map[string]Output `json:"outputs"`
	Resources        []interface{}     `json:"resources"`
	CheckResults     interface{}       `json:"check_results"`
}

type Secret struct {
	Version uint64            `json:"version"` // file format version, always 4
	Outputs map[string]Output `json:"outputs"` // variables
}

type Output struct {
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	Sensitive bool        `json:"sensitive,omitempty"`
}

type Lock struct {
	ID        string `json:"ID"`
	Operation string `json:"Operation"`
	Info      string `json:"Info"`
	Who       string `json:"Who"`
	Version   string `json:"Version"`
	Created   string `json:"Created"`
	Path      string `json:"Path"`
}
