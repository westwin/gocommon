package middleware

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// HeaderDumpConfig defines the config for http request header dump middleware.
	HeaderDumpConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper
	}
)

var (
	// DefaultHeaderDumpConfig is the default http request header dump middleware config.
	DefaultHeaderDumpConfig = HeaderDumpConfig{
		Skipper: middleware.DefaultSkipper,
	}
)

func HeaderDump() echo.MiddlewareFunc {
	return HeaderDumpWithConfig(DefaultHeaderDumpConfig)
}

// HeaderDumpWithConfig returns a http request header dump middleware with config.
func HeaderDumpWithConfig(config HeaderDumpConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultHeaderDumpConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			headers := c.Request().Header
			for k, v := range headers {
				fmt.Printf("%s:%s\n", k, v)
			}
			return next(c)
		}
	}
}
