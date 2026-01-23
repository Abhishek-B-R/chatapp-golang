package middleware

import (
	"context"
	"net/http"

	"github.com/Abhishek-B-R/chat-app-golang/internals/store"
	"github.com/Abhishek-B-R/chat-app-golang/internals/utils"
)

type ChatMiddleware struct{
	ChatMemberStore store.ChatMemberStore
}

type chatContextKey string
const ChatContextKey = chatContextKey("chat")

func SetChatMembership(r *http.Request, chatMember *store.ChatMember) *http.Request {
	ctx := context.WithValue(r.Context(), ChatContextKey, chatMember)
	return r.WithContext(ctx)
}

func GetChatMembership(r *http.Request) (*store.ChatMember, bool) {
	chatMember, ok := r.Context().Value(ChatContextKey).(*store.ChatMember)
	return chatMember, ok && chatMember != nil
}

func (cm *ChatMiddleware) RequireMembership(next http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		user, ok := GetUser(r)
		if !ok {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error":"signin to continue"})
			return
		}

		chatID, err := utils.ReadParam(r,"chatID")
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid chat ID"})
			return
		}

		isMember, err := cm.ChatMemberStore.IsMember(r.Context(), chatID, user.ID)
		if err != nil || !isMember {
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"you are not a member of this chat"})
			return
		}

		ctx := context.WithValue(r.Context(), "chatID", chatID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}