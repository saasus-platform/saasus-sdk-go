package middleware

import (
	"context"
	"net/http"

	"github.com/saasus-platform/saasus-sdk-go/ctxlib"
)

// ExtractTraceId extracts X-SaaSus-Trace-Id from request and set it to context.
func ExtractTraceId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		xSaaSusTraceId := r.Header.Get("X-SaaSus-Trace-Id")
		if xSaaSusTraceId != "" {
			ctx = context.WithValue(ctx, ctxlib.XSaaSusTraceIdKey, xSaaSusTraceId)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
