package oidc

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestIDTokenMandatoryClaims(t *testing.T) {
	tid := "example"
	sub := "johndoe"
	iss := "https://example.com"
	iat := int64(time.Now().Unix())
	exp := int64(time.Now().Add(time.Duration(2) * time.Hour).Unix())

	claims := IDTokenClaims{
		TenantID: tid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  iat,
			ExpiresAt: exp,
			Subject:   sub,
			Issuer:    iss,
		},
	}

	assert.NoError(t, claims.Valid())

	claims.TenantID = ""
	assert.Error(t, claims.Valid(), "missing tid")
	claims.TenantID = tid

	claims.StandardClaims.Issuer = ""
	assert.Error(t, claims.Valid(), "missing iss")
	claims.StandardClaims.Issuer = iss

	claims.StandardClaims.Subject = ""
	assert.Error(t, claims.Valid(), "missing sub")
	claims.StandardClaims.Subject = sub

	claims.StandardClaims.IssuedAt = 0
	assert.Error(t, claims.Valid(), "missing iat")
	claims.StandardClaims.IssuedAt = iat

	claims.StandardClaims.ExpiresAt = 0
	assert.Error(t, claims.Valid(), "missing exp")
	claims.StandardClaims.ExpiresAt = exp

	assert.NoError(t, claims.Valid())
}
