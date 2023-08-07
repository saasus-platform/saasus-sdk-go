package callback

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/modules/auth"
)

// CallbackRouteFunction is a function for callback route
func CallbackRouteFunction(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code is not provided by query parameter", http.StatusBadRequest)
		return
	}

	authClientWithResponse, err := auth.AuthWithResponse()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authFlow := authapi.GetAuthCredentialsParamsAuthFlow("tempCodeAuth")
	req := &authapi.GetAuthCredentialsParams{
		Code:         &code,
		AuthFlow:     &authFlow,
		RefreshToken: nil,
	}
	res, err := authClientWithResponse.GetAuthCredentialsWithResponse(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.JSON500 != nil {
		http.Error(w, res.JSON500.Message, http.StatusInternalServerError)
		return
	}
	if res.JSON404 != nil {
		http.Error(w, res.JSON404.Message, http.StatusUnauthorized)
		return
	}
	if res.JSON200 == nil {
		http.Error(w, "response failed", res.StatusCode())
		return
	}

	json.NewEncoder(w).Encode(res.JSON200)
}
