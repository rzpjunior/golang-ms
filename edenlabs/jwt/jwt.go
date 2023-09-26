package jwt

import (
	"github.com/golang-jwt/jwt"
)

type JWT struct {
	SignKey []byte
}

// UserClaim defines user token claim
type UserClaim struct {
	UserID      int64    `json:"user_id"`
	Permissions []string `json:"permissions"`
	ExpiresAt   int64    `json:"expire_at"`
	Timezone    string   `json:"timezone"`
	jwt.StandardClaims
}

type UserMobile struct {
	PhoneNo   string `json:"phone_no"`
	ExpiresAt int64  `json:"expire_at"`
	Timezone  string `json:"timezone"`
	jwt.StandardClaims
}

type UserCourierClaim struct {
	CourierID string `json:"courier_id"`
	SiteID    string `json:"site_id"`
	Timezone  string `json:"timezone"`
	jwt.StandardClaims
}

type UserHelperMobileClaim struct {
	UserID   string `json:"user_id"`
	SiteId   string `json:"site_id"`
	Platform string `json:"platform"`
	Timezone string `json:"timezone"`
	jwt.StandardClaims
}

func NewJWT(signKey []byte) *JWT {
	return &JWT{
		SignKey: signKey,
	}
}

func (j *JWT) Create(claims jwt.Claims) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(j.SignKey)
	return
}

func (j *JWT) Parse(token string) (tkn *jwt.Token, err error) {
	tkn, err = jwt.ParseWithClaims(token, &UserClaim{}, j.keyFunction)
	return
}
func (j *JWT) keyFunction(token *jwt.Token) (interface{}, error) {
	return []byte(j.SignKey), nil
}

func (j *JWT) Valid(token string) (valid bool, err error) {
	tkn, err := j.Parse(token)
	if err != nil {
		return
	}

	valid = tkn.Valid
	return
}

func (j *JWT) ParseMobile(token string) (tkn *jwt.Token, err error) {
	tkn, err = jwt.ParseWithClaims(token, &UserMobile{}, j.keyFunction)
	return
}

func (j *JWT) ValidMobile(token string) (valid bool, err error) {
	tkn, err := j.ParseMobile(token)
	if err != nil {
		return
	}

	valid = tkn.Valid
	return
}

func (j *JWT) ParseCourier(token string) (tkn *jwt.Token, err error) {
	tkn, err = jwt.ParseWithClaims(token, &UserCourierClaim{}, j.keyFunction)
	return
}

func (j *JWT) ValidCourier(token string) (valid bool, err error) {
	tkn, err := j.ParseCourier(token)
	if err != nil {
		return
	}

	valid = tkn.Valid
	return
}

func (j *JWT) ParseHelperMobile(token string) (tkn *jwt.Token, err error) {
	tkn, err = jwt.ParseWithClaims(token, &UserHelperMobileClaim{}, j.keyFunction)
	return
}

func (j *JWT) ValidHelperMobile(token string) (valid bool, err error) {
	tkn, err := j.ParseHelperMobile(token)
	if err != nil {
		return
	}

	valid = tkn.Valid
	return
}
