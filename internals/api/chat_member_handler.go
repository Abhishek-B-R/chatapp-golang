package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Abhishek-B-R/chat-app-golang/internals/store"
	"github.com/Abhishek-B-R/chat-app-golang/internals/utils"
)

type ChatMemberHandler struct {
	chatMemberStore store.ChatMemberStore
	logger *log.Logger
}

func NewChatMemberHandler(ChatMemberStore store.ChatMemberStore, logger *log.Logger) *ChatMemberHandler {
	return &ChatMemberHandler{
		chatMemberStore: ChatMemberStore,
		logger: logger,
	}
}

func (cmh *ChatMemberHandler) HandleAddMember(w http.ResponseWriter, r *http.Request) {
	var params struct{
		ChatID, UserID int64
		Role string
	}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		cmh.logger.Printf("ERROR: decodingAddMember: %v\n",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	err = cmh.chatMemberStore.AddMember(params.ChatID, params.UserID, params.Role)
	if err != nil {
		cmh.logger.Printf("ERROR: addMember: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to add member"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{})
}

func (cmh *ChatMemberHandler) HandleRemoveMember(w http.ResponseWriter, r *http.Request) {
	chatID, err := utils.ReadParam(r, "chatID")
	userID, err2 := utils.ReadParam(r, "userID")

	if err != nil {
		cmh.logger.Printf("ERROR: decodingRemoveMember: err: %v\n",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}
	if err2 != nil {
		cmh.logger.Printf("ERROR: decodingRemoveMember: err2: %v\n",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	err = cmh.chatMemberStore.RemoveMember(chatID, userID)
	if err != nil {
		cmh.logger.Printf("ERROR: removeMember: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to remove member"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"status":"success"})
}

func (cmh *ChatMemberHandler) HandleGetChatMembers(w http.ResponseWriter, r *http.Request){
	chatID, err := utils.ReadParam(r, "chatID")
	if err != nil {
		cmh.logger.Printf("ERROR: decodingGetChatMembers: %v\n",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	chatMembers, err := cmh.chatMemberStore.GetChatMembers(chatID)
	if err != nil {
		cmh.logger.Printf("ERROR: getChatMembers: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to retrieve chat members"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"chat_members":chatMembers})
}

func (cmh *ChatMemberHandler) HandleGetUserRole(w http.ResponseWriter, r *http.Request) {
	chatID, err := utils.ReadParam(r, "chatID")
	userID, err2 := utils.ReadParam(r, "userID")

	if err != nil || err2 != nil {
		cmh.logger.Printf("ERROR: decodingGetUserRole: %v\n",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	role, err := cmh.chatMemberStore.GetUserRole(chatID, userID)
	if err != nil {
		cmh.logger.Printf("ERROR: getUserRole: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to retrieve user role"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"role":role})
}

func (cmh *ChatMemberHandler) HandleIsMember(w http.ResponseWriter, r *http.Request) {
	var params struct{
		ChatID, UserID int64
	}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		cmh.logger.Printf("ERROR: decodingIsMember: %v\n",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	isMember, err := cmh.chatMemberStore.IsMember(params.ChatID, params.UserID)
	if err != nil {
		cmh.logger.Printf("ERROR: isMember: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to retrieve is_member"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"is_member":isMember})
}

func (cmh *ChatMemberHandler) HandleUpdateLastRead(w http.ResponseWriter, r *http.Request) {
	var params struct{
		ChatID, UserID, MessageID int64
	}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		cmh.logger.Printf("ERROR: decodingUpdateLastRead: %v\n",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	err = cmh.chatMemberStore.UpdateLastRead(params.ChatID, params.UserID, params.MessageID)
	if err != nil {
		cmh.logger.Printf("ERROR: updateLastRead: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to update last read"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{})
}

