package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/westwin/gocommon/authz"
)

type (
	// AuthzConfig defines the config for authz middleware.
	AuthzConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper
		// SessionPermsKeyName is a string to find current user's perms(default is "perms")
		SessionPermsKeyName string
	}
)

var (
	// DefaultAuthzConfig is the default authz middleware config.
	DefaultAuthzConfig = AuthzConfig{
		Skipper:             noAuthz,
		SessionPermsKeyName: "perms",
	}
)

// DefaultSkipper returns false which processes the middleware.
func noAuthz(echo.Context) bool {
	return true
}

// IsPermGrantedFunc is used to check grants
type IsPermGrantedFunc = func(needs *authz.Permission) echo.MiddlewareFunc

var HasPermFunc func() IsPermGrantedFunc = HasPerm

//HasPerm returns permission config with default authz config
func HasPerm() IsPermGrantedFunc {
	return HasPermWithConfig(DefaultAuthzConfig)
}

// HasPermWithConfig returns an HasPerm middleware with config.
func HasPermWithConfig(config AuthzConfig) IsPermGrantedFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultAuthzConfig.Skipper
	}
	if config.SessionPermsKeyName == "" {
		config.SessionPermsKeyName = "perms"
	}

	return func(needs *authz.Permission) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if config.Skipper(c) {
					return next(c)
				}
				has := c.Get(config.SessionPermsKeyName)
				if has == nil {
					return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Unauthorized, needs perm of %s", needs.ID()))
				} else if authz.PermissionsString(has.([]string)).IsGranted(needs) {
					return next(c)
				}

				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Unauthorized, needs perm of %s", needs.ID()))
			}
		}
	}
}

// PermitAll to a adapter to ignore the permission check
// HasPermFunc = PermitAll
func PermitAll() IsPermGrantedFunc {
	return func(needs *authz.Permission) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return next(c)
			}
		}
	}
}
