package dtorm

// Config is the database configuration
type Config struct {
	Type                 string `json:"type"`
	Host                 string `json:"host,omitzero"`
	Database             string `json:"database"`
	User                 string `json:"user,omitzero"`
	Password             string `json:"password,omitzero"`
	Deletable            bool   `json:"deletable,omitzero"`
	DisabledTransactions bool   `json:"disabledTransactions,omitzero"`
}
