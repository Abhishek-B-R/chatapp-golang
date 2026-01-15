package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Abhishek-B-R/chat-app-golang/internals/store"
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
		
	}
}