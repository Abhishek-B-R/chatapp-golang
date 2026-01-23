package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/Abhishek-B-R/chat-app-golang/internals/middleware"
	"github.com/Abhishek-B-R/chat-app-golang/internals/store"
	"github.com/Abhishek-B-R/chat-app-golang/internals/utils"
	"github.com/go-chi/chi"
)

type registeredUserRequest struct{
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Bio string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type UserHandler struct {
	userStore store.UserStore
	logger *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger: logger,
	}
}

func (uh *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req registeredUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: decodingCreateUser: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	if req.Username == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "username cannot be empty"})
		return
	}

	if len(req.Username) > 50 {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "username too long"})
		return
	}
	
	//lowercase the username
	req.Username = strings.ToLower(req.Username)

	err = utils.ValidateEmail(req.Email)
	if err != nil {
		uh.logger.Printf("ERROR: invalid email: %v",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid email data"})
		return
	}

	err = uh.validateRegisterRequest(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email: req.Email,
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: hashing password: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"internal server error"})
		return
	}

	err = uh.userStore.CreateUser(r.Context(), user)
	if err != nil {
		uh.logger.Printf("ERROR: registering user %v",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"internal server error"})
		return
	}
	
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user":user})
}

func (uh *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadParam(r, "userID")
	if err != nil {
		uh.logger.Printf("ERROR: decodingGetUserByID: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	user, err := uh.userStore.GetUserByID(r.Context(), userID)
	if err != nil {
		uh.logger.Printf("ERROR: getUserByID: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to get user by id"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user":user})
}

func (uh *UserHandler) HandlerGetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		uh.logger.Printf("ERROR: decodingGetUserByEmail: Invalid email field")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	user, err := uh.userStore.GetUserByEmail(r.Context(), email)
	if err != nil {
		uh.logger.Printf("ERROR: getUserByEmail: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to get user by email"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user":user})
}

func (uh *UserHandler) HandleGetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	//lowercase the username
	username = strings.ToLower(username)
	
	if username == "" {
		uh.logger.Printf("ERROR: decodingGetUserByUsername: Invalid username field")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	user, err := uh.userStore.GetUserByUsername(r.Context(), username)
	if err != nil {
		uh.logger.Printf("ERROR: getUserByUsername: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to get user by username"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user":user})
}

func (uh *UserHandler) HandleUpdateLastSeen(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, ok := middleware.GetUser(r)
	if !ok {
		uh.logger.Printf("ERROR: user not found in context")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error":"authentication required"})
		return
	}

	err := uh.userStore.UpdateLastSeen(r.Context(), authenticatedUser.ID)
	if err != nil {
		uh.logger.Printf("ERROR: updateLastSeen: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to update last seen"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{})
}

func (uh *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, ok := middleware.GetUser(r)
	if !ok {
		uh.logger.Printf("ERROR: user not found in context")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error":"authentication required"})
		return
	}

	var updateReq struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		Bio       string `json:"bio"`
		AvatarURL string `json:"avatar_url"`
	}


	err := json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		uh.logger.Printf("ERROR: decoding update request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	if updateReq.Username == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "username cannot be empty"})
		return
	}

	if len(updateReq.Username) > 50 {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "username too long"})
		return
	}

	err = utils.ValidateEmail(updateReq.Email)
	if err != nil {
		uh.logger.Printf("ERROR: invalid email: %v",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid email"})
		return
	}

	updatedUser := &store.User{
		ID:        authenticatedUser.ID,
		Email:     updateReq.Email,
		Username:  updateReq.Username, 
		Bio:       updateReq.Bio,      
		AvatarURL: updateReq.AvatarURL,
	}

	err = uh.userStore.UpdateUser(r.Context(), updatedUser)
	if err != nil {
		uh.logger.Printf("ERROR: updating user credentials : %v",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user":updatedUser})
}

func (uh *UserHandler) HandleUpdateUserPassword(w http.ResponseWriter, r *http.Request){
	authenticatedUser, ok := middleware.GetUser(r)
	if !ok {
		uh.logger.Printf("ERROR: user not found in context")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error":"authentication required"})
		return
	}
	
	var password struct{
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&password)
	if err != nil{
		uh.logger.Printf("ERROR: handleUpdateUserPassword: %v",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid credentials sent"})
		return
	}

	if password.Password == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"password cannot be empty"})
		return
	}
	if len(password.Password) < 8 {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"password too short"})
		return
	}

	err = uh.userStore.UpdateUserPassword(r.Context(), password.Password, authenticatedUser.ID)
	if err != nil {
		uh.logger.Printf("ERROR: updating user password: %v",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"msg":"updated password of this user, your tokens are been revoked, please authenticate again to continue"})
}


func (uh *UserHandler) HandleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, ok := middleware.GetUser(r)
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error":"authentication required"})
		return
	}

	user, err := uh.userStore.GetCurrentUser(r.Context(), authenticatedUser.ID)
	if err != nil {
		uh.logger.Printf("ERROR: handleGetCurrentUser: %v",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"unable to fetch user"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user":user})
}

func (h *UserHandler) validateRegisterRequest(req *registeredUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if len(req.Username) > 50 {
		return errors.New("Username cannot be greater than 50 characters")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}