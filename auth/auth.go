package auth

import (
	"net/http"
	"html/template"
	"encoding/json"
)

type Jwt struct {
	AccessToken string
	RefreshToken string
}

type User struct {
	Email string
	Password string
}

var GetTokenHandler = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

	// w.Write([]byte(token))
	w.Write([]byte("token"))
})

var RenderLoginView = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("auth/login.html")
	tmpl.Execute(w, nil)
})

var Authentication = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user := &User{Email: email, Password: password}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(user)
})
