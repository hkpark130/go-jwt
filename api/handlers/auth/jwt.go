package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
	"time"
)

type Jwt struct {
	Alg        string
	Secret_key string
}

type Header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

type Payload struct {
	Exp      time.Time `json:"exp"`
	Iat      time.Time `json:"iat"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func Hashing(payload *Payload) string {
	jwt := &Jwt{Alg: "HS256", Secret_key: "park"}
	h := hmac.New(sha256.New, []byte(jwt.Secret_key))

	jsonHeader, err := json.Marshal(Header{
		Typ: "JWT",
		Alg: jwt.Alg,
	})
	if err != nil {
		log.Fatal("json encode error: %w ", err)
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("json encode error: %w ", err)
	}

	h.Write([]byte(
		strings.Join([]string{
			base64.RawURLEncoding.EncodeToString(jsonHeader),
			base64.RawURLEncoding.EncodeToString(jsonPayload)}, ".")))

	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	token := strings.Join([]string{
		base64.RawURLEncoding.EncodeToString(jsonHeader),
		base64.RawURLEncoding.EncodeToString(jsonPayload),
		signature}, ".")

	return token
}
