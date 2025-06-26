package modulyn

type Event struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type Feature struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Enabled   bool      `json:"enabled"`
	JsonValue JsonValue `json:"jsonValue"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}

type JsonValue struct {
	Key     string   `json:"key"`
	Values  []string `json:"values"`
	Enabled bool     `json:"enabled"`
}
