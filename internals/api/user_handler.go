package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Abhishek-B-R/chat-app-golang/internals/store"
	"github.com/Abhishek-B-R/chat-app-golang/internals/utils"
)

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
	var user store.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		uh.logger.Printf("ERROR: decodingCreateUser: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	err = uh.userStore.CreateUser(&user)
	if err != nil {
		uh.logger.Printf("ERROR: createUser: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to create user"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user":user})
}

func (uh *UserHandler) HandleGetUserByID(w http.ResponseWriter, r http.Request) {
	userID, err := utils.ReadParam(&r, "userID")
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
	var email string
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		uh.logger.Printf("ERROR: decodingGetUserByEmail: %v\n", err)
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
	var username string
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		uh.logger.Printf("ERROR: decodingGetUserByUsername: %v\n", username)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"invalid request sent"})
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