package auth

// easyjson:json
type Config struct {
	DB   configDB `json:"db"`
	Port string   `json:"port"`
}

// easyjson:json
type configDB struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}
