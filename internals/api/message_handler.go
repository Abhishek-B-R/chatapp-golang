package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Abhishek-B-R/chat-app-golang/internals/middleware"
	"github.com/Abhishek-B-R/chat-app-golang/internals/store"
	"github.com/Abhishek-B-R/chat-app-golang/internals/utils"
)

type MessageHandler struct {
	store store.MessageStore
	logger *log.Logger
}

func NewMessageHandler(store store.MessageStore, logger *log.Logger) *MessageHandler {
	return &MessageHandler{
		store: store,
		logger: logger,
	}
}

func (mh *MessageHandler) HandleCreateMessage(w http.ResponseWriter, r *http.Request) {
	var msg store.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	chatID, err2 := utils.ReadParam(r, "chatID")
	if err != nil || err2 != nil {
		mh.logger.Printf("ERROR: decodingCreateMessage: %v %v\n", err, err2)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	authenticatedUser, ok := middleware.GetUser(r)
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error":"authentication required"})
		return
	}

	msg.ChatID = chatID
	msg.SenderID = &authenticatedUser.ID

	err = mh.store.CreateMessage(r.Context(), &msg)
	if err != nil {
		mh.logger.Printf("ERROR: createMessage: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to create message"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"msg":"created message"})
}

func (mh *MessageHandler) HandleGetMessage(w http.ResponseWriter, r *http.Request) {
	msgID, err := utils.ReadParam(r, "msgID")
	if err != nil {
		mh.logger.Printf("ERROR: decodingGetMessage: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	message, err := mh.store.GetMessage(r.Context(), msgID)
	if err != nil {
		mh.logger.Printf("ERROR: getMessage: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to get message"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message":message})
}

func (mh *MessageHandler) HandleGetChatMessages(w http.ResponseWriter, r *http.Request) {
	chatID, err := utils.ReadParam(r, "chatID")
	limit, err2 := utils.ReadParam(r, "limit")
	offset, err3 := utils.ReadParam(r, "offset")

	if err != nil || err2 != nil || err3 != nil {
		mh.logger.Printf("ERROR: decodingGetChatMessages: %v %v\n", err, err2)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	messages, err := mh.store.GetChatMessages(r.Context(), chatID, limit, offset)
	if err != nil {
		mh.logger.Printf("ERROR: getChatMessages: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to get chat messages"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"messages":messages})
}

func (mh *MessageHandler) HandleUpdateMessage(w http.ResponseWriter, r *http.Request) {
	msgID, err := utils.ReadParam(r, "msgID")
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid msgID"})
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mh.logger.Printf("ERROR: decodingUpdateMessage: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "content cannot be empty"})
		return
	}

	user, ok := middleware.GetUser(r)
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "authentication required"})
		return
	}

	originalMsg, err := mh.store.GetMessage(r.Context(), msgID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "message not found"})
		return
	}

	if originalMsg.SenderID == nil || user.ID != *originalMsg.SenderID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "not allowed to update this message"})
		return
	}
	originalMsg.Content = &req.Content

	if err := mh.store.UpdateMessage(r.Context(), originalMsg); err != nil {
		mh.logger.Printf("ERROR: updateMessage: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to update message"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": originalMsg})
}


func (mh *MessageHandler) HandleDeleteMessage(w http.ResponseWriter, r *http.Request){
	id, err := utils.ReadParam(r, "msgID")
	if err != nil {
		mh.logger.Printf("ERROR: decodingDeleteMessage: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	err = mh.store.DeleteMessage(r.Context(), id)
	if err != nil {
		mh.logger.Printf("ERROR: deleteMessage: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to delete message"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (mh *MessageHandler) HandleGetUnreadCount(w http.ResponseWriter, r *http.Request){
	chatID, err := utils.ReadParam(r, "chatID")
	if err != nil {
		mh.logger.Printf("ERROR: decodingGetUnreadCount: %v\n",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	authenticatedUser, ok := middleware.GetUser(r)
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error":"authentication required"})
		return
	}

	count, err := mh.store.GetUnreadCount(r.Context(), chatID, authenticatedUser.ID)
	if err != nil {
		mh.logger.Printf("ERROR: getUnreadCount: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to retrieve unread count"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"count":count})
}