package jwt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTokenSuccess(t *testing.T) {
	uc := UserClaim{
		UserID: 123,
	}

	bJWT := []byte("xxx")
	nJwt := NewJWT(bJWT)
	tokenString, err := nJwt.Create(uc)
	fmt.Println(tokenString)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsInJvbGVfY29kZSI6ImFueS1yb2xlIiwidXNlcl9wZXJtaXNzaW9uIjpbInBlcm1pc3Npb24iXX0.JEMKGuIRJjZXdU0LbGhwUqwVg1n7S9Aw2i-TMKu8YSY", tokenString)
	assert.Nil(t, err)

	token, err := nJwt.Parse(tokenString)
	assert.Nil(t, err)

	assert.Equal(t, true, token.Valid)

	valid, err := nJwt.Valid(tokenString)
	assert.Equal(t, true, valid)

	valid, err = nJwt.Valid("xxxxxxxxxxxxxxxxxxxxxxxxx")
	assert.Equal(t, false, valid)
	assert.NotNil(t, err)

}

func TestGetTokenFailed(t *testing.T) {
	uc := UserClaim{
		UserID: 0,
		// SecondaryID: "unknown",
		// ClientID:    "unknown",
		// Scopes:      nil,
	}

	bJWT := []byte("xxx")
	nJwt := NewJWT(bJWT)

	jwtString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsInNlY29uZGFyeV9pZCI6InVzZXIxMjMiLCJjbGllbnRfaWQiOiJjbGllbnQxMjMiLCJzY29wZXMiOm51bGx9.LfRA3eK6SWDxWlxPoEqHm8vN8y1Mklt1U3CimZQobSs"
	jwtGenerate, err := nJwt.Create(uc)
	assert.NotEqual(t, jwtString, jwtGenerate)
	assert.NotNil(t, jwtGenerate)
	assert.Nil(t, err)
}
