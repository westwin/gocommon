package oidc

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// IDTokenClaims extends the `jwt.StandardClaims`, but do stricter claim validation.
type IDTokenClaims struct {
	jwt.StandardClaims `json:",inline"`
	TenantID           string   `json:"tid,omitempty"`
	UID                string   `json:"preferred_username,omitempty"`
	Permissions        []string `json:"perms,omitempty"`
	Roles              []string `json:"roles,omitempty"`
	Groups             []string `json:"groups,omitempty"`
	Email              string   `json:"email,omitempty"`
	EmailVerified      bool     `json:"email_verified,omitempty"`
	Phone              string   `json:"phone_number,omitempty"`
	PhoneVerified      bool     `json:"phone_number_verified,omitempty"`
	AuthorizingParty   string   `json:"azp,omitempty"`
	AccessTokenHash    string   `json:"at_hash,omitempty"`
}

// Valid implements the interface of jwt.Claims.Valid()
// There is no accounting for clock skew.
// The jwt.StandardClaims will consider as valid claim if any of the above claims are not in the token.
// But for IDTokenClaims they are mandatory claims.
func (c IDTokenClaims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	return c.mandatory()
}

func (c IDTokenClaims) mandatory() error {
	if c.TenantID == "" ||
		c.Subject == "" ||
		c.ExpiresAt == 0 ||
		c.IssuedAt == 0 ||
		c.Issuer == "" {
		return fmt.Errorf("token is invalid, missing mandatory claims")
	}

	return nil
}
