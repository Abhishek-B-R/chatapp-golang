package api

import (
	"encoding/json"
	"log"
	"net/http"

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
	if err != nil {
		mh.logger.Printf("ERROR: decodingCreateMessage: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	err = mh.store.CreateMessage(&msg)
	if err != nil {
		mh.logger.Printf("ERROR: createMessage: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to create message"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{})
}

func (mh *MessageHandler) HandleGetMessage(w http.ResponseWriter, r *http.Request) {
	msgID, err := utils.ReadParam(r, "msgID")
	if err != nil {
		mh.logger.Printf("ERROR: decodingGetMessage: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	message, err := mh.store.GetMessage(msgID)
	if err != nil {
		mh.logger.Printf("ERROR: getMessage: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to get message"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message":message})
}

func (mh *MessageHandler) HandleGetChatMessages(w http.ResponseWriter, r *http.Request) {
	var params struct {
		ChatID int64
		Limit int64
		Offset int64
	}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		mh.logger.Printf("ERROR: decodingGetChatMessages: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	messages, err := mh.store.GetChatMessages(params.ChatID, params.Limit, params.Offset)
	if err != nil {
		mh.logger.Printf("ERROR: getChatMessages: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to get chat messages"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"messages":messages})
}

func (mh *MessageHandler) HandleUpdateMessage(w http.ResponseWriter, r *http.Request){
	var msg store.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		mh.logger.Printf("ERROR: decodingUpdateMessages: %v\n",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	err = mh.store.UpdateMessage(&msg)
	if err != nil {
		mh.logger.Printf("ERROR: updateMessage: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to update message"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{})
}

func (mh *MessageHandler) HandleDeleteMessage(w http.ResponseWriter, r *http.Request){
	id, err := utils.ReadParam(r, "msgID")
	if err != nil {
		mh.logger.Printf("ERROR: decodingDeleteMessage: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	err = mh.store.DeleteMessage(id)
	if err != nil {
		mh.logger.Printf("ERROR: deleteMessage: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to delete message"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (mh *MessageHandler) HandleGetUnreadCount(w http.ResponseWriter, r *http.Request){
	var params struct{
		ChatID int64
		UserID int64
	}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		mh.logger.Printf("ERROR: decodingGetUnreadCount: %v\n",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	count, err := mh.store.GetUnreadCount(params.ChatID, params.UserID)
	if err != nil {
		mh.logger.Printf("ERROR: getUnreadCount: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to retrieve unread count"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"count":count})
}