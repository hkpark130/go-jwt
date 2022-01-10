package auth

import (
	"testing"
	"time"
)

const TEST_TOKEN = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOiIyMDIxLTExLTE3VDIwOjM0OjU4LjY1MTM4NzIzN1oiLCJpYXQiOiIyMDIxLTExLTE3VDIwOjM0OjU4LjY1MTM4NzIzN1oiLCJlbWFpbCI6InRlc3RAdGVzdC5jb20ifQ.BDIt0KLqjYxbZaCMxzK0sb5uZDBBtKfPvrTbLijCcKk"

func TestHashing(t *testing.T) {
	// 期限切れの検証で引っかかるので、2021年11月に設定
	token, err := IssueAccessToken(&Payload{
		Exp:   time.Date(2021, 11, 17, 20, 34, 58, 651387237, time.UTC),
		Iat:   time.Date(2021, 11, 17, 20, 34, 58, 651387237, time.UTC),
		Email: "test@test.com",
	})
	if err != nil {
		t.Error(err)
	}

	expectedToken := TEST_TOKEN

	if token != expectedToken {
		t.Error("Something wrong with the hashing process.")
		t.Error("test token:", token)
	}
}

func TestDecoding(t *testing.T) {
	token := TEST_TOKEN
	claim, err := Decode(token)
	if err != nil {
		t.Error(err)
	}

	if claim.Email != "test@test.com" || claim.Password != "test" {
		t.Error("Something wrong with the decoding process.")
		t.Error("test claim:", claim)
	}
}

func TestVerifyToken(t *testing.T) {
	token := TEST_TOKEN
	IsTokenVerified(token)
}
