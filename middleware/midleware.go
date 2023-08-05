package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Anti-Pattern-Inc/saasus-sdk-go/generated/authapi"
	"github.com/Anti-Pattern-Inc/saasus-sdk-go/modules/auth"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

type IDTokenGetter interface {
	GetIDToken(*http.Request) string
}

func Authenticate(w http.ResponseWriter, r *http.Request, idToken string) (*authapi.UserInfo, error) {
	if idToken == "" {
		return nil, errors.New("invalid Authorization header")
	}

	authClient, err := auth.AuthWithResponse()
	if err != nil {
		return nil, err
	}

	res, err := authClient.GetUserInfoWithResponse(context.Background(), &authapi.GetUserInfoParams{Token: idToken})
	if err != nil {
		return nil, err
	}
	if res.JSON500 != nil {
		return nil, fmt.Errorf("data = %s, message = %s, type = %s", res.JSON500.Data, res.JSON500.Message, res.JSON500.Type)
	}
	if res.JSON401 != nil {
		http.Redirect(w, r, os.Getenv("SAASUS_LOGIN_URL"), http.StatusUnauthorized)
		return nil, fmt.Errorf("data = %s, message = %s, type = %s", res.JSON401.Data, res.JSON401.Message, res.JSON401.Type)
	}

	return res.JSON200, nil
}

func AuthMiddleware(next http.Handler, getter IDTokenGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := Authenticate(w, r, getter.GetIDToken(r)); err != nil {
			http.Error(w, "Unauthorized "+err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthMiddlewareGin(getter IDTokenGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, err := Authenticate(c.Writer, c.Request, getter.GetIDToken(c.Request))
		if err != nil {
			http.Error(c.Writer, "Unauthorized "+err.Error(), http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("userInfo", userInfo)
		c.Next()
	}
}

func AuthMiddlewareEcho(getter IDTokenGetter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userInfo, err := Authenticate(c.Response().Writer, c.Request(), getter.GetIDToken(c.Request()))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "Unauthorized "+err.Error())
			}

			c.Set("userInfo", userInfo)
			return next(c)
		}
	}
}
