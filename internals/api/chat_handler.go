package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Abhishek-B-R/chat-app-golang/internals/store"
	"github.com/Abhishek-B-R/chat-app-golang/internals/utils"
)

type ChatHandler struct {
	chatStore store.ChatStore
	logger *log.Logger
}

func NewChatHandler(chatStore store.ChatStore, logger *log.Logger) *ChatHandler {
	return &ChatHandler{
		chatStore: chatStore,
		logger: logger,
	}
}

func (ch *ChatHandler) HandleCreateChat(w http.ResponseWriter, r *http.Request) {
	var chat store.Chat
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		ch.logger.Printf("ERROR: decodingCreateChat: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
		return
	}

	createdChat, err := ch.chatStore.CreateChat(&chat)
	if err != nil {
		ch.logger.Printf("ERROR: createChat: %v\n",err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"Failed to create chat"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"chat":createdChat})
}

func (ch *ChatHandler) HandleGetChatByID(w http.ResponseWriter, r *http.Request) {
	chatID, err := utils.ReadIDParam(r)
	if err != nil {
		ch.logger.Printf("ERROR: readIDParam: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid chat id"})
		return
	}

	chat, err := ch.chatStore.GetChatByID(int(chatID))
	if err != nil {
		ch.logger.Printf("ERROR: GetWorkoutByID: %v\n", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error":"internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"chat":chat})
}

func (ch *ChatHandler) HandleUpdateChat(w http.ResponseWriter, r *http.Request) {
	var chat store.Chat
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		ch.logger.Printf("ERROR: decodingUpdateChat: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid request sent"})
	}

	err = ch.chatStore.UpdateChat(&chat)
	if err != nil {
		ch.logger.Printf("ERROR: updateChat: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error":"failed to update chat"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"chat":chat})
}

func (ch *ChatHandler) HandleDeleteChat(w http.ResponseWriter, r *http.Request) {
	chatID, err := utils.ReadIDParam(r)
	if err != nil {
		ch.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid chat delete id"})
		return
	}

	err = ch.chatStore.DeleteChat(int(chatID))
	if err == sql.ErrNoRows {
		ch.logger.Printf("ERROR: deleteChat: %v\n",err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"No such chat found in db"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}