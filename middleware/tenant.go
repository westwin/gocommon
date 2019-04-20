package middleware

import (
	"net/http"

	"github.com/labstack/echo"
)

type (
	// TenantResolverConfig defines the config for RenantResolver middleware.
	TenantResolverConfig struct {
		// the HTTP Header name which represents tenant id.
		// Note: the HTTP Header takes precedence over `TenantQueryParam`
		TenantHeader string
		// the param of the query string which represents tenant id
		TenantQueryParam string
		// default tenant id in case we can not resolve the tenant from query string
		DefaultTenant string
		// the context key of tenant id
		SessionTenantKeyName string
	}
)

var (
	// DefaultTenantResolverConfig is the default TenantResolver middleware config.
	DefaultTenantResolverConfig = TenantResolverConfig{
		TenantHeader:         "X-TID",
		TenantQueryParam:     "tid",
		SessionTenantKeyName: "tid",
		DefaultTenant:        "example"}
)

// TenantResolver returns a root level (before router) middleware which resolves
// the tenant id from query string and set it as a content.
//
// Usage `Echo#Pre(TenantResolver())`
func TenantResolver() echo.MiddlewareFunc {
	return TenantResolverWithConfig(DefaultTenantResolverConfig)
}

// TenantResolverWithConfig returns a TenantResolver middleware with config.
// See `TenantResolver()`.
func TenantResolverWithConfig(config TenantResolverConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tid := c.Request().Header.Get(config.TenantHeader)
			if tid == "" {
				tid = c.QueryParam(config.TenantQueryParam)
			}
			if tid == "" {
				tid = config.DefaultTenant
			}

			if tid != "" {
				// authz test
				//c.Set("perms", []string{"*.*.*"})
				c.Set(config.SessionTenantKeyName, tid)
				return next(c)
			}

			return echo.NewHTTPError(http.StatusBadRequest, "missing tenant id")
		}
	}
}
