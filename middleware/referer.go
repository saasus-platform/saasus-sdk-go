package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/saasus-platform/saasus-sdk-go/ctxlib"
)

// ExtractReferer extracts referer from request and set it to context.
func ExtractReferer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ref := r.Referer()
		if ref != "" {
			ctx = context.WithValue(ctx, ctxlib.RefererKey, ref)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ExtractRefererGin extracts referer from request and set it to context.
func ExtractRefererGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ref := c.Request.Referer()
		if ref != "" {
			ctx := context.WithValue(c.Request.Context(), ctxlib.RefererKey, ref)
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}

// ExtractRefererEcho extracts referer from request and set it to context.
func ExtractRefererEcho() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ref := c.Request().Referer()
			if ref != "" {
				ctx := context.WithValue(c.Request().Context(), ctxlib.RefererKey, ref)
				c.SetRequest(c.Request().WithContext(ctx))
			}

			return next(c)
		}
	}
}
