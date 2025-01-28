package internal

type NewLogin struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

var users = map[string]NewLogin{}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
