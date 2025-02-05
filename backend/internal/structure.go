package internal

import "log"

type NewLogin struct {
	Username       string
	HashedPassword string
	SessionToken   string
	CSRFToken      string
	IsAdmin        bool
}

func initUsers() map[string]*NewLogin {
	adminHashedPassword, err := hashPassword("adminpass")
	if err != nil {
		log.Fatalf("Ошибка хеширования пароля администратора: %v", err)
	}

	userHashedPassword, err := hashPassword("userpass")
	if err != nil {
		log.Fatalf("Ошибка хеширования пароля пользователя: %v", err)
	}

	return map[string]*NewLogin{
		"admin": {Username: "admin", HashedPassword: adminHashedPassword, IsAdmin: true},
		"user":  {Username: "user", HashedPassword: userHashedPassword, IsAdmin: false},
	}
}

var users = initUsers()

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
