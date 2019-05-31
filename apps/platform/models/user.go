package models

// easyjson:json
type UserProfile struct {
	UID      int64  `json:"uid"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
