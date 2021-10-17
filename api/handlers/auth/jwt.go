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

func Decode(token string) Payload {
	jwt := &Jwt{Alg: "HS256", Secret_key: os.Getenv("SECRET_KEY")}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		log.Fatal("Invalid JWT Structure")
	}
	header, _ := base64.RawURLEncoding.DecodeString(parts[0])
	payload, _ := base64.RawURLEncoding.DecodeString(parts[1])
	signature := parts[2]
	// JSON decode payload
	var pldat Payload
	if err := json.Unmarshal(payload, &pldat); err != nil {
		log.Fatal(err.Error())
	}
	// JSON decode header
	var headdat Header
	if err := json.Unmarshal(header, &headdat); err != nil {
		log.Fatal(err.Error())
	}
	// Extract and parse expiration date from payload
	layout := "2006-01-02 15:04:05"
	exp := pldat.Exp.Format(layout)
	expParsed, err := time.ParseInLocation(layout, exp, time.Now().Location())
	if err != nil {
		log.Fatal(err)
	}
	// Check how old the JWT is.  Return an error if it is expired
	now := time.Now()
	if now.After(expParsed) {
		log.Fatal("Expired JWT")
	}
	// This probably should be one of the first checks, preceeding the date check.  If the signature of the JWT doesn't match there is likely fuckery afoot
	ha := hmac256(string(parts[0])+"."+string(parts[1]), jwt.Secret_key)
	if ha != string(signature) {
		log.Fatal("Invalid JWT signature")
	}

	return pldat
}
