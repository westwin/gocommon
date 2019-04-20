package oidc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"io"
	"strings"

	jose "gopkg.in/square/go-jose.v2"
)

// auth reqeust scopes
const (
	ScopeOfflineAccess     = "offline_access"
	ScopeAddress           = "address"
	ScopeOpenID            = "openid"                    // Request a id_token
	ScopeGroups            = "groups"                    // Request a user group
	ScopeEmail             = "email"                     // Request user email
	ScopePhone             = "phone"                     // Request user phone
	ScopeProfile           = "profile"                   // Request user profile
	ScopeCrossClientPrefix = "audience:proxy:client_id:" // Proxy auth
)

// auth request response type
const (
	ResponseTypeCode    = "code"     // "Regular" flow
	ResponseTypeToken   = "token"    // Implicit flow for frontend apps.
	ResponseTypeIDToken = "id_token" // ID Token in url fragment
)

// auth reqeust grant type
const (
	GrantTypeAuthorizationCode  = "authorization_code"
	GrantTypeRefreshToken       = "refresh_token"
	GrantTypeClientCredential   = "client_credentials"
	GrantTypePasswordCredential = "password"
	GrantTypeImplicit           = "implicit"
)

// SupportedResponseTypes will return in discovery endpoint
var SupportedResponseTypes = []string{ResponseTypeToken, ResponseTypeIDToken, ResponseTypeCode}

// Scopes represents additional data requested by the clients about the end user.
type Scopes struct {
	// The client has requested a refresh token from the server.
	OfflineAccess bool

	// The client has requested group information about the end user.
	Groups bool

	// The client has requested openid info about the end user.
	Openid bool
}

// DefaultScopes will be used if no scopes is requested
// it doesn't contain scopeCrossClientPrefix, usually this scopeCrossClientPrefix must be explicitly specified by request
var DefaultScopes = []string{
	ScopeOfflineAccess,
	ScopeOpenID,
	ScopeGroups,
	ScopeEmail,
	ScopeProfile,
}

// ScopeDescriptions describe scopes in consent page
var ScopeDescriptions = map[string]string{
	"offline_access": "Have offline access",
	"profile":        "View basic profile information",
	"email":          "View your email address",
}

// ParseScopes will flag scopes according to given scopes string
func ParseScopes(scopes []string) Scopes {
	var s Scopes
	for _, scope := range scopes {
		switch scope {
		case ScopeOfflineAccess:
			s.OfflineAccess = true
		case ScopeGroups:
			s.Groups = true
		case ScopeOpenID:
			s.Openid = true
		}
	}
	return s
}

//CompareScopes compares requested scopes with the specified scopes, returns authorized and unauthorized scopes
func CompareScopes(requestedScopes, scopes []string) (authorizedScopes, unauthorizedScopes []string) {
	for _, s := range requestedScopes {
		contains := func() bool {
			for _, scope := range scopes {
				if s == scope {
					return true
				}
			}
			return false
		}()
		if contains {
			authorizedScopes = append(authorizedScopes, s)
		} else {
			unauthorizedScopes = append(unauthorizedScopes, s)
		}
	}
	return authorizedScopes, unauthorizedScopes
}

// NeedRefreshToken check whether auth request need refresh_token
func NeedRefreshToken(scopes []string) bool {
	for _, scope := range scopes {
		if scope == ScopeOfflineAccess {
			return true
		}
	}
	return false
}

// SignatureAlgorithm is for determine the signature algorithm for a JWT.
func SignatureAlgorithm(jwk *jose.JSONWebKey) (alg jose.SignatureAlgorithm, err error) {
	if jwk.Key == nil {
		return alg, errors.New("no signing key")
	}
	switch key := jwk.Key.(type) {
	case *rsa.PrivateKey:
		// Because OIDC mandates that we support RS256, we always return that
		// value. In the future, we might want to make this configurable on a
		// per client basis. For example allowing PS256 or ECDSA variants.
		//
		// See https://toi/oidc/issues/692
		return jose.RS256, nil
	case *ecdsa.PrivateKey:
		// We don't actually support ECDSA keys yet, but they're tested for
		// in case we want to in the future.
		//
		// These values are prescribed depending on the ECDSA key type. We
		// can't return different values.
		switch key.Params() {
		case elliptic.P256().Params():
			return jose.ES256, nil
		case elliptic.P384().Params():
			return jose.ES384, nil
		case elliptic.P521().Params():
			return jose.ES512, nil
		default:
			return alg, errors.New("unsupported ecdsa curve")
		}
	default:
		return alg, fmt.Errorf("unsupported signing key type %T", key)
	}
}

// SignPayload will sign a payload using given key and algorithm
func SignPayload(key *jose.JSONWebKey, alg jose.SignatureAlgorithm, payload []byte) (jws string, err error) {
	signingKey := jose.SigningKey{Key: key, Algorithm: alg}

	signer, err := jose.NewSigner(signingKey, &jose.SignerOptions{})
	if err != nil {
		return "", fmt.Errorf("new signier: %v", err)
	}
	signature, err := signer.Sign(payload)
	if err != nil {
		return "", fmt.Errorf("signing payload: %v", err)
	}
	return signature.CompactSerialize()
}

// The hash algorithm for the at_hash is detemrined by the signing
// algorithm used for the id_token. From the spec:
//
//    ...the hash algorithm used is the hash algorithm used in the alg Header
//    Parameter of the ID Token's JOSE Header. For instance, if the alg is RS256,
//    hash the access_token value with SHA-256
//
// https://openid.net/specs/openid-connect-core-1_0.html#ImplicitIDToken
var hashForSigAlg = map[jose.SignatureAlgorithm]func() hash.Hash{
	jose.RS256: sha256.New,
	jose.RS384: sha512.New384,
	jose.RS512: sha512.New,
	jose.ES256: sha256.New,
	jose.ES384: sha512.New384,
	jose.ES512: sha512.New,
}

// AccessTokenHash Compute an at_hash from a raw access token and a signature algorithm
//
// See: https://openid.net/specs/openid-connect-core-1_0.html#ImplicitIDToken
func AccessTokenHash(alg jose.SignatureAlgorithm, accessToken string) (string, error) {
	newHash, ok := hashForSigAlg[alg]
	if !ok {
		return "", fmt.Errorf("unsupported signature algorithm: %s", alg)
	}

	hash := newHash()
	if _, err := io.WriteString(hash, accessToken); err != nil {
		return "", fmt.Errorf("computing hash: %v", err)
	}
	sum := hash.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(sum[:len(sum)/2]), nil
}

// ParseCrossClientScope extract the peerID from specified scope
func ParseCrossClientScope(scope string) (peerID string, ok bool) {
	if ok = strings.HasPrefix(scope, ScopeCrossClientPrefix); ok {
		peerID = scope[len(ScopeCrossClientPrefix):]
	}
	return
}
