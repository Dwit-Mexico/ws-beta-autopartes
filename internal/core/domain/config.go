package domain

type DSNSourceConfig struct {
	Name string
	DSN  string
}

type DatabaseColumnDefinition struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
