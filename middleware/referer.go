package middleware

import (
	"context"
	"net/http"

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
		xSaaSusReferer := r.Header.Get("X-SaaSus-Referer")
		if xSaaSusReferer != "" {
			ctx = context.WithValue(ctx, ctxlib.XSaaSusRefererKey, xSaaSusReferer)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
