package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

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

	err = uh.userStore.CreateUser(user)
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

	user, err := uh.userStore.GetUserByID(userID)
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

	user, err := uh.userStore.GetUserByEmail(email)
	if err != nil {
		uh.logger.Printf("ERROR: getUserByEmail: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to get user by email"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user":user})
}

func (uh *UserHandler) HandleGetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		uh.logger.Printf("ERROR: decodingGetUserByUsername: Invalid username field")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	user, err := uh.userStore.GetUserByUsername(username)
	if err != nil {
		uh.logger.Printf("ERROR: getUserByUsername: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to get user by email"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user":user})
}

func (uh *UserHandler) HandleUpdateLastSeen(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadParam(r, "userID")
	if err != nil {
		uh.logger.Printf("ERROR: decodingUpdateLastSeen: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"invalid request sent"})
		return
	}

	err = uh.userStore.UpdateLastSeen(userID)
	if err != nil {
		uh.logger.Printf("ERROR: updateLastSeen: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to update last seen"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{})
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