package middleware

import (
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/westwin/gocommo/oidc"
)

type (
	// JWTConfig defines the config for JWT middleware.
	JWTConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// BeforeFunc defines a function which is executed just before the middleware.
		BeforeFunc middleware.BeforeFunc

		// SuccessHandler defines a function which is executed for a valid token.
		SuccessHandler JWTSuccessHandler

		// ErrorHandler defines a function which is executed for an invalid token.
		// It may be used to define a custom JWT error.
		ErrorHandler JWTErrorHandler

		// Signing method, used to check token signing method.
		// Optional. Default value RS256.
		SigningMethod string

		// Context key to store user information from the token into context.
		// Optional. Default value "user".
		ContextKey string

		// Claims are extendable claims data defining token content.
		// Optional. Default value common.IDTokenClaims
		Claims jwt.Claims

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		TokenLookup string

		// AuthScheme to be used in the Authorization header.
		// Optional. Default value "Bearer".
		AuthScheme string

		// KeyLookup lookups the key per tid/kid from backed storage
		KeyLookup KeyLookupFunc

		keyFunc jwt.Keyfunc
	}

	// JWTSuccessHandler defines a function which is executed for a valid token.
	JWTSuccessHandler func(echo.Context)

	// JWTErrorHandler defines a function which is executed for an invalid token.
	JWTErrorHandler func(error) error

	jwtExtractor func(echo.Context) (string, error)
)

// Algorithms
const (
	AlgorithmRS256 = "RS256"
)

// Errors
var (
	ErrJWTMissing = echo.NewHTTPError(http.StatusBadRequest, "missing or malformed access token")
)

var (
	// DefaultJWTConfig is the default JWT auth middleware config.
	DefaultJWTConfig = JWTConfig{
		Skipper:       middleware.DefaultSkipper,
		SigningMethod: AlgorithmRS256,
		ContextKey:    "user",
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		AuthScheme:    "Bearer",
		Claims:        oidc.IDTokenClaims{},
	}
)

// KeyLookupFunc is the function to look up the key
type KeyLookupFunc func(tid, kid string) (interface{}, error)

// JWT returns a JSON Web Token (JWT) auth middleware.
//
// For valid token, it sets the user in context and calls next handler.
// For invalid token, it returns "401 - Unauthorized" error.
// For missing token, it returns "400 - Bad Request" error.
//
// See: https://jwt.io/introduction
// See `JWTConfig.TokenLookup`
func JWT(keyLookup KeyLookupFunc) echo.MiddlewareFunc {
	c := DefaultJWTConfig
	c.KeyLookup = keyLookup
	return JWTWithConfig(c)
}

// JWTWithConfig returns a JWT auth middleware with config.
// See: `JWT()`.
func JWTWithConfig(config JWTConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultJWTConfig.Skipper
	}
	if config.KeyLookup == nil {
		panic("echo: jwt middleware requires signing key lookup func")
	}
	if config.SigningMethod == "" {
		config.SigningMethod = DefaultJWTConfig.SigningMethod
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultJWTConfig.ContextKey
	}
	if config.Claims == nil {
		config.Claims = DefaultJWTConfig.Claims
	}
	if config.TokenLookup == "" {
		config.TokenLookup = DefaultJWTConfig.TokenLookup
	}
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultJWTConfig.AuthScheme
	}
	config.keyFunc = func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != config.SigningMethod {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}

		// kid MUST be formatted like <tid>/<kid>
		kid := t.Header["kid"]
		if kid == nil {
			return nil, fmt.Errorf("missing kid")
		}

		if pair := strings.Split(kid.(string), "/"); len(pair) == 2 {
			return config.KeyLookup(pair[0], pair[1])
		}
		return nil, fmt.Errorf("invalid kid")
	}

	// Initialize
	parts := strings.Split(config.TokenLookup, ":")
	extractor := jwtFromHeader(parts[1], config.AuthScheme)
	switch parts[0] {
	case "query":
		extractor = jwtFromQuery(parts[1])
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			if config.BeforeFunc != nil {
				config.BeforeFunc(c)
			}

			auth, err := extractor(c)
			if err != nil {
				if config.ErrorHandler != nil {
					return config.ErrorHandler(err)
				}
				return err
			}
			token := new(jwt.Token)
			if _, ok := config.Claims.(oidc.IDTokenClaims); ok {
				token, err = jwt.ParseWithClaims(auth, &oidc.IDTokenClaims{}, config.keyFunc)
			}

			if err == nil && token.Valid {
				reqTenantID := c.Get("tid").(string)
				claims := token.Claims.(*oidc.IDTokenClaims)
				if reqTenantID != claims.TenantID { // token binding tid is not the requesting tenant
					return echo.NewHTTPError(http.StatusBadRequest,
						fmt.Sprintf("requesting tid:%s, but token_tid:%s", reqTenantID, claims.TenantID))
				}

				// Store user information from token into context.
				//c.Set(config.ContextKey, token)
				c.Set("perms", claims.Permissions)
				if config.SuccessHandler != nil {
					config.SuccessHandler(c)
				}
				return next(c)
			}
			if config.ErrorHandler != nil {
				return config.ErrorHandler(err)
			}
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  "invalid or expired token",
				Internal: err,
			}
		}
	}
}

// jwtFromHeader returns a `jwtExtractor` that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}

// jwtFromQuery returns a `jwtExtractor` that extracts token from the query string.
func jwtFromQuery(param string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		token := c.QueryParam(param)
		if token == "" {
			return "", ErrJWTMissing
		}
		return token, nil
	}
}
