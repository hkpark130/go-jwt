package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"
	// "strings"
	"github.com/gin-gonic/gin"
)

func TestTokenPathHandler(t *testing.T) {
	router := gin.Default()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/token", nil)

	router.ServeHTTP(res, req)
	// GetTokenHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatal("Not 200 Status / ", res.Code)
	}

	data, _ := ioutil.ReadAll(res.Body)

	if e := string(data); e != "token" {
		t.Errorf("Wrong Response %s", e)
	}
}

// func TestAuthenticationHandler(t *testing.T) {
// 	res := httptest.NewRecorder()
// 	req := httptest.NewRequest("POST", "/api/login", 
// 		strings.NewReader("email=hkpark@kddi.com&password=1234") )
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	Authentication(res, req)

// 	if res.Code != http.StatusOK {
// 		t.Fatal("Not 200 Status / ", res.Code)
// 	}

// 	var user User
// 	json.NewDecoder(res.Body).Decode(&user)

// 	if e := user.Email; e != "hkpark@kddi.com" {
// 		t.Errorf("Email doesn't match %s != %s", "hkpark@kddi.com", e)
// 	}
	
// }

// req.Header.Add("Authorization", "auth_token=\"~\"")
