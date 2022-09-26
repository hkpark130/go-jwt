// Package auth is JWT 토큰 인증관련 패키지
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

// Jwt 는 토큰 해쉬화 때 사용됨
type Jwt struct {
	Alg       string
	SecretKey string
}

// Header 은 토큰 타입과 해쉬 알고리즘 정보를 가지고 있음
type Header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

// Payload 는 유저 데이터를 가지고 있음
type Payload struct {
	Exp        time.Time `json:"exp"`
	Iat        time.Time `json:"iat"`
	Email      string    `json:"email"`
	Permission string    `json:"permission"`
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

// ModifyForRefreshToken 은 RefreshToken 에 넣을 Payload 을 다시 만듬
func ModifyForRefreshToken(payload *Payload) *Payload {
	modifiedPayload := &Payload{}
	*modifiedPayload = *payload
	modifiedPayload.Email = ""
	modifiedPayload.Exp = payload.Exp.Add(time.Hour * 24 * 1)
	modifiedPayload.Permission = payload.Permission

	return modifiedPayload
}

// IssueToken 은 파라미터로 받은 Payload 을 기반으로 토큰을 만듬
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

// ReissueToken 은 파라미터로 받은 Access Token 의 payload 을 기반으로 토큰을 재발행
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

	newPayload := CreatePayload(payload.Email, payload.Permission)
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

func CreatePayload(email string, permission string) *Payload {
	payload := &Payload{
		Exp:        time.Now().Add(time.Second * time.Duration(3600)), //3600 = 1H
		Iat:        time.Now(),
		Permission: permission,
		Email:      email}

	return payload
}

// isExpired 은 파라미터로 받은 Payload 의 유효기간을 확인
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

// Decode 은 파라미터로 받은 token 을 Payload 에 디코딩함
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

// IsTokenVerified 은 파라미터로 받은 token 의 유효성을 확인
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
