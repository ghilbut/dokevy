package terraform

type Secret struct {
	FormatVersion string `json:"format_version"`
	Values        Values `json:"values"`
}

type Values struct {
	Outputs map[string]Output `json:"outputs"`
}

type Output struct {
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	Sensitive bool        `json:"sensitive"`
}
