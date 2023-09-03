package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/saasus-platform/saasus-sdk-go/ctxlib"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/modules/auth"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

// IDTokenGetter is an interface for getting ID token from request.
type IDTokenGetter interface {
	// GetIDToken returns ID token from request.
	//
	// For example, you can get it from query parameters or headers.
	GetIDToken(*http.Request) string
}

// IDTokenGetterFromAuthHeader is an implementation of IDTokenGetter.
type IdTokenGetterFromAuthHeader struct{}

// GetIDToken is a function that returns an IDToken from Authorization header.
func (*IdTokenGetterFromAuthHeader) GetIDToken(r *http.Request) string {
	// Authorization: Bearer <id_token>
	strs := strings.Split(r.Header.Get("Authorization"), " ")
	if len(strs) == 2 && strings.ToLower(strs[0]) == "bearer" {
		return strs[1]
	}
	return ""
}

// Authenticate is a function that authenticates the request and returns a UserInfo.
//
// ref: https://docs.saasus.io/reference/getuserinfo
func Authenticate(w http.ResponseWriter, r *http.Request, idToken string) (*authapi.UserInfo, error) {
	if idToken == "" {
		return nil, errors.New("invalid Authorization header")
	}

	authClient, err := auth.AuthWithResponse()
	if err != nil {
		return nil, err
	}

	res, err := authClient.GetUserInfoWithResponse(r.Context(), &authapi.GetUserInfoParams{Token: idToken})
	if err != nil {
		return nil, err
	}
	if res.JSON500 != nil {
		return nil, fmt.Errorf("data = %s, message = %s, type = %s", res.JSON500.Data, res.JSON500.Message, res.JSON500.Type)
	}
	if res.JSON401 != nil {
		return nil, fmt.Errorf("data = %s, message = %s, type = %s", res.JSON401.Data, res.JSON401.Message, res.JSON401.Type)
	}

	return res.JSON200, nil
}

// AuthMiddleware is a middleware for authentication by http.Handler.
func AuthMiddleware(next http.Handler, getter IDTokenGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo, err := Authenticate(w, r, getter.GetIDToken(r))
		if err != nil {
			http.Error(w, "Unauthorized "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxlib.UserInfoKey, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthMiddlewareGin is a middleware for authentication by gin.HandlerFunc.
func AuthMiddlewareGin(getter IDTokenGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, err := Authenticate(c.Writer, c.Request, getter.GetIDToken(c.Request))
		if err != nil {
			http.Error(c.Writer, "Unauthorized "+err.Error(), http.StatusUnauthorized)
			return
		}

		c.Set(string(ctxlib.UserInfoKey), userInfo)
		c.Next()
	}
}

// AuthMiddlewareEcho is a middleware for authentication by echo.MiddlewareFunc.
func AuthMiddlewareEcho(getter IDTokenGetter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userInfo, err := Authenticate(c.Response().Writer, c.Request(), getter.GetIDToken(c.Request()))
			if err != nil {
				http.Error(c.Response().Writer, "Unauthorized "+err.Error(), http.StatusUnauthorized)
				return nil
			}

			c.Set(string(ctxlib.UserInfoKey), userInfo)
			return next(c)
		}
	}
}
