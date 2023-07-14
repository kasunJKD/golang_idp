package handlers

import (
	//"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"idp-service/gateway/config"
	//"idp-service/gateway/utils"

	pbtoken "idp-service/protos/token"
	//pblogin "idp-service/protos/login"
)

// func handleAuthCodeAuth(w http.ResponseWriter, r *http.Request) {
// 	queryParams := r.URL.Query()
// 	clientID := queryParams.Get("client_id")

// 	switch clientID {
// 	case "":
// 		//utils.ShowError(w, r, 400, "Bad Request", "client_id is required")
// 		http.Error(w, "method not allowed", http.StatusBadRequest)
// 	case clientID:
// 		//check clientId valid
// 		//redirect to consent page with params

// 		//utils.PresentAuthScreen(w, r, config.AuthCode)
// 	default:
// 		//utils.ShowError(w, r, 401, "Unauthorized", "Invalid client_id")
// 		http.Error(w, "method not allowed", http.StatusBadRequest)
// 	}
// }

// func handleAuthSignInAuth(w http.ResponseWriter, r *http.Request) {
// 	queryParams := r.URL.Query()
// 	clientID := queryParams.Get("client_id")

// 	//create new login repo
// 	loginrepo := client.ConnectLogin(os.Getenv("LOGIN_REPO_HOST"), os.Getenv("LOGIN_REPO_PORT"))
// 	//create response from google response
// 	rr := &pblogin.ClientReq{
// 		ClientId: clientID,
// 	}
// 	_, err := loginrepo.ValidateClientId(context.Background(), rr)

// 	if err != nil {
// 		http.Error(w, "method not allowed", http.StatusBadRequest)
// 		return
// 	}

// 	switch clientID {
// 	case "":
// 		//utils.ShowError(w, r, 400, "Bad Request", "client_id is required")
// 		http.Error(w, "method not allowed", http.StatusBadRequest)
// 	case clientID:
// 		utils.PresentSignInScreen(w, r, config.AuthCode)
// 	default:
// 		//utils.ShowError(w, r, 401, "Unauthorized", "Invalid client_id")
// 		http.Error(w, "method not allowed", http.StatusBadRequest)
// 	}
// }

// handleAuthCodeToken checks for the existence of all parameters detailed in Section 4.1.3 of RFC 6749 (https://tools.ietf.org/html/rfc6749#section-4.1.3).
// If not present, an HTTP 400 response is sent.
// Else, a new token is generated, added to the store, and returned to the user in a JSON response.
func (h *HandlerFunc) handleAuthCodeToken(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if params["client_id"] == "" || params["grant_type"] == "" || params["code"] == "" {
		http.Error(w, "client_id, grant_type=authorization_code, code and redirect_uri are required", http.StatusBadRequest)
		return
	}
	//create request
	rr := &pbtoken.TokenRequest{
		Code: params["code"],
		RedirectURI: params["redirect_uri"],
	}
	token, err := h.t.Client.NewAuthCodeToken(r.Context() , rr)
	log.Println(token)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	jsonBytes, err := json.Marshal(token)
	//w.Write(jsonBytes)
	fmt.Fprintln(w, string(jsonBytes))
}

// Refer RFC 6749 Section 6 (https://tools.ietf.org/html/rfc6749#section-6)
func (h *HandlerFunc) handleAuthCodeRefresh(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// If found, invalidate previously issued token
	//create request
	rr := &pbtoken.RefreshTokenRequest{
		RefreshToken: params["refresh_token"],
		InvalidateIfFound: true,
	}
	res, _ := h.t.Client.AuthCodeRefreshTokenExists(r.Context(), rr)
	if res.Value {
		rq := &pbtoken.RefreshTokenRequest{
			RefreshToken: params["refresh_token"],
		}
		token, err := h.t.Client.NewAuthCodeRefreshToken(r.Context(), rq)
		if err != nil {
			http.Error(w, "Token generation failed. Please try again.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		jsonBytes, err := json.Marshal(token)

		fmt.Fprintln(w, string(jsonBytes))
	} else {
		http.Error(w, "expired or invalid refresh token", http.StatusInternalServerError)
	}
}
