package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"idp-service/gateway/utils"

	"idp-service/gateway/config"
	pblogin "idp-service/protos/login"
	pbtoken "idp-service/protos/token"
)

	type Payload struct {
		Signinbtn string `json:"signinbtn"`
		Email string `json:"email"`
		Password string `json:"password"`
		State string `json:"state"`
		Client_id string `json:"client_id"`
		RedirectURI string `json:"redirect_uri"`
		Response string `json:"response"`
		Flow string `json:"flow"`
		Project_id string `json:"project_id"`
		User_id string `json:"user_id"`
		Otp_code int32 `json:"otp_code"`
	}

	type Response struct {
		Url string `json:"url"`
		UserId string `json:"userId"`
		AccessToken string `json:"accessToken"`
		OtpCode string `json:"otp_code"`
	}

//Routes the request to a AuthorizationHandler based on the request_type
// func HandleAuth(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	params := r.URL.Query()

// 	// Perform empty checks on the following parameters:
// 	// - response_type
// 	// - client_id
// 	if params.Get("response_type") == "" || params.Get("client_id") == "" {
// 		http.Error(w, "method not allowed", http.StatusBadRequest)
// 		return
// 	}

// 	switch r.URL.Query().Get("response_type") {
// 	case "code":
// 		handleAuthCodeAuth(w, r)
// 	default:
// 		http.Error(w, "method not allowed", http.StatusBadRequest)
// 		//utils.ShowError(w, r, http.StatusBadRequest, "Authorization Flow Error", "Unknown response_type: "+r.URL.Query().Get("response_type"))
// 	}
// }

// func HandleOauthSignIn(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	params := r.URL.Query()

// 	// Perform empty checks on the following parameters:
// 	// - response_type
// 	// - client_id
// 	if params.Get("response_type") == "" || params.Get("client_id") == "" {
// 		http.Error(w, "method not allowed", http.StatusBadRequest)
// 		return
// 	}

// 	switch r.URL.Query().Get("response_type") {
// 	case "code":
// 		handleAuthSignInAuth(w, r)
// 	default:
// 		http.Error(w, "method not allowed", http.StatusBadRequest)
// 		//utils.ShowError(w, r, http.StatusBadRequest, "Authorization Flow Error", "Unknown response_type: "+r.URL.Query().Get("response_type"))
// 	}
// }

func (h*HandlerFunc) HandleOauthSignInResponse(w http.ResponseWriter, r *http.Request) {	
	contents, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var payload Payload
	err = json.Unmarshal(contents, &payload)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}


    signinbtn := payload.Signinbtn
	email := payload.Email
	password := payload.Password
	state := payload.State
	clientId := payload.Client_id
	projectId := payload.Project_id
	redirectURI, err := url.QueryUnescape(payload.RedirectURI)


	log.Println(payload)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//need to fix back to GETENV
	authorizeUrl := "http://" + h.env.Server.Host + ":44201/membership/authorize"
	//clientId := r.FormValue("client_id")

	if err != nil {
		//utils.ShowError(w, r, http.StatusBadRequest, "Bad Request", "Invalid redirect_uri")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	
	if signinbtn == "signin" {
		var res Response
		rr := &pblogin.Request{
			Email: email,
			Password:  password,
		}
		//calling SignUpWithAuthProvider
		response, err := h.l.Client.PasswordSignIn(context.Background(), rr)
		
		if err != nil {
			log.Println(err)
			http.Error(w, "User not exists", http.StatusBadRequest)
			return
		}

		userID := response.Users.UserId

		if response.Users.OtpEnabled == 1 {
			authorizeUrl = "http://" + h.env.Server.Host + ":44201/membership/otpLogin/oauth"
			// return the otpCode in response for now, as aws services still not connected
			res.OtpCode = response.Users.OtpCode
		} else {
			//add user id to cache
			rq := &pbtoken.User{
				UserId: userID,
			}
			_, err = h.t.Client.AddUserIdAuthCodeFlow(context.Background(), rq)
			if err != nil {
				log.Println(err)
			}
		}
		authorizeUrl += "?response_type=code&client_id=" + clientId + "&redirect_uri=" + redirectURI + "&state=" + state + "&project_id=" + projectId + "&userId=" + userID

		res.Url = authorizeUrl
		res.UserId = userID
		res.AccessToken = response.OauthAccessToken
		resp, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(resp)
		
	} 

	if signinbtn == "P99IDP" {
		var res Response
		rr := &pblogin.Request{
			Email: email,
			Password:  password,
		}
		//calling SignUpWithAuthProvider
		response, err := h.l.Client.PasswordSignIn(context.Background(), rr)
		
		if err != nil {
			log.Println(err)
			http.Error(w, "User not exists", http.StatusBadRequest)
			return
		}

		userID := response.Users.UserId

		if response.Users.OtpEnabled == 1 {
			authorizeUrl = "http://" + h.env.Server.Host + ":44201/membership/otpLogin/oauth"
			// return the otpCode in response for now, as aws services still not connected
			res.OtpCode = response.Users.OtpCode
		} else {
			//add user id to cache
			rq := &pbtoken.User{
				UserId: userID,
			}
			_, err = h.t.Client.AddUserIdAuthCodeFlow(context.Background(), rq)
			if err != nil {
				log.Println(err)
			}
			authorizeUrl = "http://" + h.env.Server.Host + ":44203/projects"
		}
		//authorizeUrl = redirectURI
		authorizeUrl += "?response_type=access&client_id=p99signInAuthentication" + "&redirect_uri=" + redirectURI + "&state=" + state + "&userId=" + userID + "&access=" + response.OauthAccessToken

		res.Url = authorizeUrl
		res.UserId = userID
		res.AccessToken = response.OauthAccessToken
		resp, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(resp)
	}
	
	
	//http.Redirect(w, r, authorizeUrl, http.StatusTemporaryRedirect)
}

func (h*HandlerFunc) HandleOauthOtpLoginResponse(w http.ResponseWriter, r *http.Request) {
	contents, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var payload Payload
	err = json.Unmarshal(contents, &payload)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}


	signinbtn := payload.Signinbtn
	state := payload.State
	clientId := payload.Client_id
	projectId := payload.Project_id
	redirectURI, err := url.QueryUnescape(payload.RedirectURI)
	userId := payload.User_id
	otpCode := payload.Otp_code

	log.Println(payload)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//need to fix back to GETENV
	authorizeUrl := "http://" + h.env.Server.Host + ":44201/membership/authorize"
	//clientId := r.FormValue("client_id")

	if err != nil {
		//utils.ShowError(w, r, http.StatusBadRequest, "Bad Request", "Invalid redirect_uri")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if signinbtn == "signin" {
		var res Response
		rr := &pblogin.OtpLoginRequest{
			UserId: userId,
			OtpCode:  otpCode,
		}
		//calling SignUpWithAuthProvider
		response, err := h.l.Client.OtpLogin(context.Background(), rr)

		if err != nil {
			log.Println(err)
			http.Error(w, "otp code is invalid", http.StatusBadRequest)
			return
		}

		//add user id to cache
		rq := &pbtoken.User{
			UserId: userId,
		}
		_, err = h.t.Client.AddUserIdAuthCodeFlow(context.Background(), rq)
		if err != nil {
			log.Println(err)
		}
		authorizeUrl += "?response_type=code&client_id=" + clientId + "&redirect_uri=" + redirectURI + "&state=" + state + "&project_id=" + projectId

		res.Url = authorizeUrl
		res.UserId = userId
		res.AccessToken = response.OauthAccessToken
		resp, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(resp)

	}

	if signinbtn == "P99IDP" {
		var res Response
		rr := &pblogin.OtpLoginRequest{
			UserId: userId,
			OtpCode:  otpCode,
		}
		//calling SignUpWithAuthProvider
		response, err := h.l.Client.OtpLogin(context.Background(), rr)

		if err != nil {
			log.Println(err)
			http.Error(w, "otp code is invalid", http.StatusBadRequest)
			return
		}

		//add user id to cache
		rq := &pbtoken.User{
			UserId: userId,
		}
		_, err = h.t.Client.AddUserIdAuthCodeFlow(context.Background(), rq)
		if err != nil {
			log.Println(err)
		}
		authorizeUrl = "http://" + h.env.Server.Host + ":44203/projects"
		//authorizeUrl = redirectURI
		authorizeUrl += "?response_type=access&client_id=p99signInAuthentication" + "&redirect_uri=" + redirectURI + "&state=" + state + "&userId=" + userId + "&access=" + response.OauthAccessToken

		res.Url = authorizeUrl
		res.UserId = userId
		res.AccessToken = response.OauthAccessToken
		resp, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(resp)
	}


	//http.Redirect(w, r, authorizeUrl, http.StatusTemporaryRedirect)
}

// Invoked by the Authorization Grant screen when the user accepts the authorization request.
// Extracts the redirect_uri from the JSON body, attaches an authorization grant to it,
// and redirects the user-agent to that URI.
func (h*HandlerFunc) HandleResponse(w http.ResponseWriter, r *http.Request) {
	contents, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var payload Payload
	err = json.Unmarshal(contents, &payload)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	flow, err := strconv.Atoi(payload.Flow)
	if err != nil {
		//utils.ShowError(w, r, 400, "OAuth 2.0 Flow Error", "Unrecognized flow")
		//return
		http.Error(w, "OAuth 2.0 Flow Error", http.StatusNotFound)
		return
	}

	response := payload.Response
	state := payload.State
	projectId := payload.Project_id
	clientId := payload.Client_id
	redirectURI, err := url.QueryUnescape(payload.RedirectURI)
	if err != nil {
		//utils.ShowError(w, r, http.StatusBadRequest, "Bad Request", "Invalid redirect_uri")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if response == "ACCEPT" {
		switch flow {
		case config.AuthCode:
			rq := &pbtoken.TokenRequest{
				RedirectURI: redirectURI,
			}
			req, _ := h.t.Client.NewAuthCodeGrant(r.Context(), rq)
			redirectURI += "?code=" + req.Value + "&state=" + state + "&project_id=" + projectId + "&client_id=" + clientId
		}
	} else if response == "CANCEL" {
		redirectURI += "?error=access_denied"
	}

	url, err := json.Marshal(redirectURI)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(url)
	//http.Redirect(w, r, redirectURI, http.StatusSeeOther)
}

// Redirects the request to the appropriate flowHandler by checking the 'grant_type' parameter.
// Refer RFC 6749 Section 4.1.3 (https://tools.ietf.org/html/rfc6749#section-4.1.3)
// Accepts only POST requests with application/x-www-form-urlencoded body.
func (h*HandlerFunc) HandleToken(w http.ResponseWriter, r *http.Request) {
	log.Println("hit success HandleToken")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//utils.ShowJSONError(w, r, 500, "An error occurred while processing your request")
		http.Error(w, "An error occurred while processing your request", http.StatusInternalServerError)
		return
	}

	params, err := utils.ParseParams(string(body))
	if err != nil {
		log.Println(err)
		//utils.ShowJSONError(w, r, 400, "Expected parameters not found.")
		http.Error(w, "Expected parameters not found.", http.StatusNotFound)
		return
	}

	if params["client_id"] == "" && params["client_secret"] == "" {
		clientID, clientSecret := utils.ParseBasicAuthHeader(r.Header.Get("Authorization"))
		params["client_id"] = clientID
		params["client_secret"] = clientSecret
	}

	switch params["grant_type"] {
	case "authorization_code":
		h.handleAuthCodeToken(w, r, params)
	case "refresh_token":
		if len(params["refresh_token"]) != 72 {
			http.Error(w, "refresh_token missing or invalid", http.StatusNotFound)
			return
		}
		
		if strings.HasPrefix(params["refresh_token"], config.AuthCodeFlowID) {
			h.handleAuthCodeRefresh(w, r, params)
		}
	default:
		http.Error(w, "grant_type absent or invalid", http.StatusNotFound)
		
	}
}