package internal

import (
	"errors"
	"fmt"
	"net/http"
)

func Authorize(r *http.Request) error {
	username := r.FormValue("username")
	fmt.Println("Authorize Username:", username)

	user, ok := users[username]
	if !ok {
		fmt.Println("User not found")
		return errors.New("user not found")
	}
	fmt.Printf("Authorize User Data: %+v\n", user)

	st, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Session token cookie missing or error:", err)
		return errors.New("session_token cookie not found")
	}
	fmt.Println("Session Token from Cookie:", st.Value)
	if st.Value != user.SessionToken {
		fmt.Println("Session token mismatch")
		return errors.New("invalid session token")
	}

	csrf := r.Header.Get("X-CSRF-Token")
	fmt.Println("CSRF Token from Header:", csrf)
	if csrf != user.CSRFToken {
		fmt.Println("CSRF token mismatch")
		return errors.New("invalid CSRF token")
	}

	fmt.Println("Authorization successful")
	return nil
}
