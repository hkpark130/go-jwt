// Package auth is JWT トークン認証関連パッケージ
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"golang/jwt/api/repository"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

// Jwt はトークンハッシュ化の際に使われる
type Jwt struct {
	Alg       string
	SecretKey string
}

// Header はトークンタイプとハッシュアルゴリズムの情報を持っている
type Header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

// Payload はユーザーデータを持っている
type Payload struct {
	Exp   time.Time `json:"exp"`
	Iat   time.Time `json:"iat"`
	Email string    `json:"email"`
}

var (
	errInvalidJwt   = errors.New("Invalid JWT Structure")
	errExpiredToken = errors.New("Expired JWT Token")
)

func hmac256(message, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// ModifyForRefreshToken はRefreshToken に入れるPayload を作り直す
func ModifyForRefreshToken(payload *Payload) *Payload {
	modifiedPayload := &Payload{}
	*modifiedPayload = *payload
	modifiedPayload.Email = ""
	modifiedPayload.Exp = payload.Exp.Add(time.Hour * 24 * 7)

	return modifiedPayload
}

// IssueToken は受け取ったPayload をもとにトークンを作る
func IssueToken(payload *Payload) (string, error) {
	jwt := &Jwt{Alg: "HS256", SecretKey: os.Getenv("SECRET_KEY")}

	jsonHeader, err := json.Marshal(Header{
		Typ: "JWT",
		Alg: jwt.Alg,
	})
	if err != nil {
		log.Panicln("json encode error: %w ", err)
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Panicln("json encode error: %w ", err)
	}

	msg := strings.Join([]string{
		base64.RawURLEncoding.EncodeToString(jsonHeader),
		base64.RawURLEncoding.EncodeToString(jsonPayload)}, ".")

	signature := hmac256(msg, jwt.SecretKey)

	token := strings.Join([]string{
		base64.RawURLEncoding.EncodeToString(jsonHeader),
		base64.RawURLEncoding.EncodeToString(jsonPayload),
		signature}, ".")

	return token, err
}

// ReissueToken は受け取った Access Token の payload をもとにトークンを再発行する
func ReissueToken(payload Payload, jwtUserRepository *repository.JwtUserRepository) (string, string, error) {
	refreshToken, err := jwtUserRepository.GetRefreshToken(payload.Email)
	if err != nil {
		if err == redis.Nil {
			return "", "", errExpiredToken
		} else {
			log.Println(err)
			return "", "", err
		}
	}

	pldat, err := Decode(refreshToken)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	if !isExpired(pldat) {
		return "", "", errExpiredToken
	}

	newPayload := CreatePayload(payload.Email)
	accessToken, err := IssueToken(newPayload)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	modifiedPayload := ModifyForRefreshToken(newPayload)
	refreshToken, err = IssueToken(modifiedPayload)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	err = jwtUserRepository.SetRefreshToken(payload.Email, refreshToken, modifiedPayload.Exp)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func CreatePayload(email string) *Payload {
	payload := &Payload{
		Exp:   time.Now().Add(time.Second * time.Duration(3600)), //3600 = 1H
		Iat:   time.Now(),
		Email: email}

	return payload
}

// isExpired は受け取ったPayload の有効期限を確認する
func isExpired(pldat Payload) bool {
	layout := "2006-01-02 15:04:05"
	exp := pldat.Exp.Format(layout)
	expParsed, err := time.ParseInLocation(layout, exp, time.Now().Location())
	if err != nil {
		log.Println("failed: %w", err)
		return false
	}

	now := time.Now()
	if now.After(expParsed) {
		log.Println("expired JWT")
		return false
	}
	return true
}

func parseJWT(token string) ([]string, error) {
	if regexp.MustCompile(`^[\w-]+\.[\w-]+\.[\w-]+$`).MatchString(token) {
		parts := strings.Split(token, ".")
		return parts, nil
	}
	log.Println(errInvalidJwt)
	return nil, errInvalidJwt
}

// Decode は受け取ったtoken をPayload にデコーディングする
func Decode(token string) (Payload, error) {
	parts, err := parseJWT(token)
	if err != nil {
		log.Println(err)
	}

	decodedPayload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Println(err)
	}

	var pldat Payload
	if err := json.Unmarshal(decodedPayload, &pldat); err != nil {
		log.Println(err.Error())
	}

	return pldat, err
}

// IsTokenVerified は受け取ったtoken の有効性を確認する
func IsTokenVerified(token string) (bool, error) {
	jwt := &Jwt{Alg: "HS256", SecretKey: os.Getenv("SECRET_KEY")} //id+pw
	parts, err := parseJWT(token)
	if err != nil {
		log.Println(err)
		return false, err
	}

	pldat, err := Decode(token)
	if err != nil {
		log.Println(err)
		return false, err
	}

	if !isExpired(pldat) {
		return false, errExpiredToken
	}

	ha := hmac256(strings.Join([]string{parts[0], parts[1]}, "."), jwt.SecretKey)
	if ha != string(parts[2]) {
		log.Println("invalid JWT signature")
		return false, errInvalidJwt
	}

	return true, nil
}
