package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"
	"strings"
	"golang/jwt/app"
	"golang/jwt/auth"
	"encoding/json"
)

func TestTokenPathHandler(t *testing.T) {
	router := app.SetupRouter()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/token", nil)

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatal("Not 200 Status / ", res.Code)
	}

	data, _ := ioutil.ReadAll(res.Body)

	if e := string(data); e != "token" {
		t.Errorf("Wrong Response %s", e)
	}
}

func TestAuthenticationHandler(t *testing.T) {
	router := app.SetupRouter()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/login", 
		strings.NewReader("email=hkpark@kddi.com&password=1234") )
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatal("Not 200 Status / ", res.Code)
	}

	var user auth.User
	json.NewDecoder(res.Body).Decode(&user)

	if e := user.Email; e != "hkpark@kddi.com" {
		t.Errorf("Email doesn't match %s != %s", "hkpark@kddi.com", e)
	}	
}

// func TestRenderLoginViewHandler(t *testing.T) {
// 	router := app.SetupRouter()

// 	res := httptest.NewRecorder()
// 	req := httptest.NewRequest("GET", "/login", nil)

// 	router.ServeHTTP(res, req)
// }
