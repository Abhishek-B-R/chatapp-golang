package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Abhishek-B-R/chat-app-golang/internals/store"
	"github.com/Abhishek-B-R/chat-app-golang/internals/utils"
)

type TokenHandler struct{
	tokenStore store.TokenStore
	userStore store.UserStore
	logger *log.Logger
}

type createTokenRequest struct{
	UserName string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler{
	return &TokenHandler{tokenStore: tokenStore, userStore: userStore, logger: logger}
}

func (th *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		th.logger.Printf("ERROR: createTokenRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request payload"})
		return
	}

	user, err := th.userStore.GetUserByUsername(r.Context(), req.UserName)
	if err != nil || user == nil {
		th.logger.Printf("ERROR: GetUserByUsername: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error":"internal server error"})
		return
	}

	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		th.logger.Printf("ERROR: PasswordHash.Matches %v",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"internal server error"})
		return
	}

	if !passwordsDoMatch {
		fmt.Println()
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error":"invalid credentials"})
		return
	}

	token, err := th.tokenStore.CreateNewToken(r.Context(), user.ID, 24*time.Hour)
	if err != nil {
		th.logger.Printf("ERROR: Creating token %v",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token":token})
}