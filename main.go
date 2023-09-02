// main.go

package main

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	users      = map[string]string{"admin": "password"} // User credentials (in-memory for simplicity)
	sessionKey = "authenticated"                        // Session key to track user authentication
	mu         sync.Mutex                               // Mutex for concurrent access to session data
)

func main() {
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.ListenAndServe(":8080", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if isValidUser(username, password) {
			// Set session flag to indicate authentication
			setSession(w)
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	// Display login form
	fmt.Fprintln(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Login Page</title>
		</head>
		<body>
			<h1>Login Page</h1>
			<form method="POST">
				<label for="username">Username:</label>
				<input type="text" id="username" name="username" required><br>
				<label for="password">Password:</label>
				<input type="password" id="password" name="password" required><br>
				<button type="submit">Login</button>
			</form>
		</body>
		</html>
	`)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fmt.Fprintln(w, "Welcome to the Dashboard!")
}

func isValidUser(username, password string) bool {
	// Check user credentials (in-memory for simplicity)
	storedPassword, ok := users[username]
	return ok && storedPassword == password
}

func setSession(w http.ResponseWriter) {
	mu.Lock()
	defer mu.Unlock()
	http.SetCookie(w, &http.Cookie{
		Name:  sessionKey,
		Value: "true",
	})
}

func isAuthenticated(r *http.Request) bool {
	mu.Lock()
	defer mu.Unlock()
	cookie, err := r.Cookie(sessionKey)
	return err == nil && cookie.Value == "true"
}
