package models

// easyjson:json
type App struct {
	AID           int64  `json:"uid"`
	Link          string `json:"link"`
	URL           string `json:"url"`
	Name          string `json:"name"`
	About         string `json:"about"`
	Image         string `json:"image"`
	Installations int64  `json:"installs"`
	Category      string `json:"category"`
}

// easyjson:json
type Apps []App
