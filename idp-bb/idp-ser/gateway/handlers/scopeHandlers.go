package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	pbtoken "idp-service/protos/token"
	pbuser "idp-service/protos/user"
)

func (h*HandlerFunc) HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	params := r.URL.Query()
	accesstoken := params.Get("access_token")

	if params.Get("access_token") == "" {
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}
	
	rq := &pbtoken.TokenRequest{
		Token: accesstoken,
	}
	//validate accesstoken
	b, _ := h.t.Client.VerifyAuthCodeToken(r.Context(), rq)
	if b.Value == false {
		http.Error(w, "token expired", http.StatusBadRequest)
	}

	rqq := &pbtoken.User{
		AccessToken: accesstoken,
	}
	id, err := h.t.Client.GetUserIdfromAccesstoken(r.Context(), rqq)
	log.Println(id.Value + " this is id from idp gateway")

	if err != nil {
		http.Error(w, "token expired", http.StatusBadRequest)
	}
	
	//create response from google response
	rr := &pbuser.Request{
		UserId: id.Value,
	}
	w.Header().Set("Content-Type", "application/json")
	res, err := h.u.Client.GetUserInfoById(context.Background(), rr)
	log.Println(res)
	if err != nil {
		http.Error(w, "User id not found", http.StatusInternalServerError)
	}
	//w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(res)

}