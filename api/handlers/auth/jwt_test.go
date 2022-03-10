package auth

import (
	"testing"
	"time"
)

const TEST_ACCESS_TOKEN = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOiIyMDIxLTExLTE3VDIwOjM0OjU4LjY1MTM4NzIzN1oiLCJpYXQiOiIyMDIxLTExLTE3VDIwOjM0OjU4LjY1MTM4NzIzN1oiLCJlbWFpbCI6InRlc3RAdGVzdC5jb20ifQ.BDIt0KLqjYxbZaCMxzK0sb5uZDBBtKfPvrTbLijCcKk"
const TEST_REFRESH_TOKEN = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOiIyMDIxLTExLTI0VDIwOjM0OjU4LjY1MTM4NzIzN1oiLCJpYXQiOiIyMDIxLTExLTE3VDIwOjM0OjU4LjY1MTM4NzIzN1oiLCJlbWFpbCI6IiJ9.ZzRYvGk7svCuRs41GlHJEXLGMtWFDc0BCm4XT2pWMnk"

func TestHashing(t *testing.T) {
	// 期限切れの検証で引っかかるので、2021年11月に設定
	payload := &Payload{
		Exp:   time.Date(2021, 11, 17, 20, 34, 58, 651387237, time.UTC),
		Iat:   time.Date(2021, 11, 17, 20, 34, 58, 651387237, time.UTC),
		Email: "test@test.com",
	}

	accessToken, err := IssueToken(payload)
	if err != nil {
		t.Error(err)
	}

	modifiedPayload := ModifyForRefreshToken(payload)
	refreshToken, err := IssueToken(modifiedPayload)
	if err != nil {
		t.Error(err)
	}

	if accessToken != TEST_ACCESS_TOKEN || refreshToken != TEST_REFRESH_TOKEN {
		t.Error("Something wrong with the hashing process.")
		t.Error("test access token:", accessToken)
		t.Error("test refresh token:", refreshToken)
	}
}

func TestDecoding(t *testing.T) {
	token := TEST_ACCESS_TOKEN
	claim, err := Decode(token)
	if err != nil {
		t.Error(err)
	}

	if claim.Email != "test@test.com" {
		t.Error("Something wrong with the decoding process.")
		t.Error("test claim:", claim)
	}
}

func TestVerifyToken(t *testing.T) {
	token := TEST_ACCESS_TOKEN
	IsTokenVerified(token)
}
