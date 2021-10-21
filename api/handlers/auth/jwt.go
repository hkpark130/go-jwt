package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
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

const JWTLEN = 3

func hmac256(message, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func Hashing(payload *Payload) string {
	jwt := &Jwt{Alg: "HS256", Secret_key: os.Getenv("SECRET_KEY")}

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

	msg := strings.Join([]string{
		base64.RawURLEncoding.EncodeToString(jsonHeader),
		base64.RawURLEncoding.EncodeToString(jsonPayload)}, ".")

	signature := hmac256(msg, jwt.Secret_key)

	token := strings.Join([]string{
		base64.RawURLEncoding.EncodeToString(jsonHeader),
		base64.RawURLEncoding.EncodeToString(jsonPayload),
		signature}, ".")

	return token
}

func isExpired(pldat Payload) {
	layout := "2006-01-02 15:04:05"
	exp := pldat.Exp.Format(layout)
	expParsed, err := time.ParseInLocation(layout, exp, time.Now().Location())
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	if now.After(expParsed) {
		log.Fatal("Expired JWT")
	}
}

func parseJWT(token string) (string, string, string) {
	parts := strings.Split(token, ".")
	if len(parts) != JWTLEN {
		log.Fatal("Invalid JWT Structure")
	}

	return parts[0], parts[1], parts[2]
}

func jsonUnmarshal(jsonBytes []byte, decodedData interface{}) {
	if err := json.Unmarshal(jsonBytes, &decodedData); err != nil {
		log.Fatal(err.Error())
	}
}

func Decode(token string) Payload {
	jwt := &Jwt{Alg: "HS256", Secret_key: os.Getenv("SECRET_KEY")}

	header, payload, signature := parseJWT(token)
	decodedPayload, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		log.Fatal(err)
	}
	var pldat Payload
	jsonUnmarshal(decodedPayload, &pldat)

	isExpired(pldat)

	ha := hmac256(strings.Join([]string{header, payload}, "."), jwt.Secret_key)
	if ha != string(signature) {
		log.Fatal("Invalid JWT signature")
	}

	return pldat
}
