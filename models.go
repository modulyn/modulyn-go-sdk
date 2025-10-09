package modulyn

type Event struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type Feature struct {
	ID        string    `json:"id"`
	Label     string    `json:"label"`
	Enabled   bool      `json:"enabled"`
	JsonValue JsonValue `json:"jsonValue"`
}

type JsonValue struct {
	Key     string   `json:"key"`
	Values  []string `json:"values"`
	Enabled bool     `json:"enabled"`
}
