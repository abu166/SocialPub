package internal

type NewLogin struct {
	Username       string
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

var users = map[string]NewLogin{}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
