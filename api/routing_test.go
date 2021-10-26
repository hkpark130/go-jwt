package app

import (
	"golang/jwt/api/adapter"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: 各Requestに対するtest code作成

func TestTokenPathHandler(t *testing.T) {
	db, err := adapter.Init()
	if err != nil {
		log.Printf("Failed to connect to Database %s ", err)
	}

	router := SetupRouter(db)
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
