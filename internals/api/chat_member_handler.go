package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Abhishek-B-R/chat-app-golang/internals/middleware"
	"github.com/Abhishek-B-R/chat-app-golang/internals/store"
	"github.com/Abhishek-B-R/chat-app-golang/internals/utils"
)

type ChatMemberHandler struct {
	chatMemberStore store.ChatMemberStore
	messageStore store.MessageStore
	logger *log.Logger
}

func NewChatMemberHandler(ChatMemberStore store.ChatMemberStore, MessageStore store.MessageStore, logger *log.Logger) *ChatMemberHandler {
	return &ChatMemberHandler{
		chatMemberStore: ChatMemberStore,
		messageStore: MessageStore,
		logger: logger,
	}
}

func (cmh *ChatMemberHandler) HandleAddMember(w http.ResponseWriter, r *http.Request) {
	var params struct{
		UserID int64
		Role string
	}
	err := json.NewDecoder(r.Body).Decode(&params)

	chatID, err2 := utils.ReadParam(r, "chatID")
	if err != nil || err2 != nil {
		cmh.logger.Printf("ERROR: decodingAddMember: %v %v\n",err, err2)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	err = cmh.chatMemberStore.AddMember(r.Context(), chatID, params.UserID, params.Role)
	if err != nil {
		cmh.logger.Printf("ERROR: addMember: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to add member"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"msg":"added user"})
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

	err = cmh.chatMemberStore.RemoveMember(r.Context(), chatID, userID)
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

	chatMembers, err := cmh.chatMemberStore.GetChatMembers(r.Context(), chatID)
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
		cmh.logger.Printf("ERROR: decodingGetUserRole: %v %v\n",err, err2)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	role, err := cmh.chatMemberStore.GetUserRole(r.Context(), chatID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "user is not a member of this chat"})
			return
		}
		
		cmh.logger.Printf("ERROR: getUserRole: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to retrieve user role"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"role":role})
}

func (cmh *ChatMemberHandler) HandleIsMember(w http.ResponseWriter, r *http.Request) {
	chatID, err := utils.ReadParam(r, "chatID")
	userID, err2 := utils.ReadParam(r, "userID")

	if err != nil || err2 != nil {
		cmh.logger.Printf("ERROR: decodingIsMember: %v %v\n",err, err2)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	isMember, err := cmh.chatMemberStore.IsMember(r.Context(), chatID, userID)
	if err != nil {
		cmh.logger.Printf("ERROR: isMember: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to retrieve is_member"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"is_member":isMember})
}

func (cmh *ChatMemberHandler) HandleUpdateLastRead(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r)
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "authentication required"})
		return
	}

	chatID, err := utils.ReadParam(r, "chatID")
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid chat ID"})
		return
	}

	var req struct {
		MessageID int64 `json:"message_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request"})
		return
	}

	msg, err := cmh.messageStore.GetMessage(r.Context(), req.MessageID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "message not found"})
		return
	}

	if msg.ChatID != chatID {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "message does not belong to this chat"})
		return
	}

	err = cmh.chatMemberStore.UpdateLastRead(r.Context(), chatID, user.ID, req.MessageID)
	if err != nil {
		cmh.logger.Printf("ERROR: updating last read: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to update"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"msg": "marked as read",
		"last_read_message_id": req.MessageID,
	})
}

func (cmh *ChatMemberHandler) HandleMuteChat(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r)
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "authentication required"})
		return
	}

	chatID, err := utils.ReadParam(r, "chatID")
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"no chatID field"})
		return
	}

	err = cmh.chatMemberStore.MuteChat(r.Context(), user.ID, chatID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"msg":"muted chat successfully"})
}

func (cmh *ChatMemberHandler) HandleUnMuteChat(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r)
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "authentication required"})
		return
	}

	chatID, err := utils.ReadParam(r, "chatID")
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"no chatID field"})
		return
	}

	err = cmh.chatMemberStore.MuteChat(r.Context(), user.ID, chatID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"msg":"unmuted chat successfully"})
}
