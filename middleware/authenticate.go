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
)

// IDTokenGetter is an interface for getting ID token from request.
type IDTokenGetter interface {
	// GetIDToken returns ID token from request.
	//
	// For example, you can get it from query parameters or headers.
	GetIDToken(*http.Request) string
}

// IdTokenGetterFromAuthHeader is an implementation of IDTokenGetter.
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
func Authenticate(ctx context.Context, idToken string) (*authapi.UserInfo, error) {
	if idToken == "" {
		return nil, errors.New("invalid Authorization header")
	}

	authClient, err := auth.AuthWithResponse()
	if err != nil {
		return nil, err
	}

	res, err := authClient.GetUserInfoWithResponse(ctx, &authapi.GetUserInfoParams{Token: idToken})
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
		userInfo, err := Authenticate(r.Context(), getter.GetIDToken(r))
		if err != nil {
			http.Error(w, "Unauthorized "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxlib.UserInfoKey, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
