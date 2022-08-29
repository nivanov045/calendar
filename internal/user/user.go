package user

type User struct {
	Info CustomUserInfo `json:"info"`         // info about user
	ID   string         `json:"id,omitempty"` //id
}

type CustomUserInfo struct {
	Name string `json:"name"` // user's name
}
