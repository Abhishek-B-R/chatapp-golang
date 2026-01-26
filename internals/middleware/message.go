package middleware

import (
	"context"
	"net/http"

	"github.com/Abhishek-B-R/chat-app-golang/internals/store"
	"github.com/Abhishek-B-R/chat-app-golang/internals/utils"
)

type MessageMiddleware struct{
	MessageStore store.MessageStore
	ChatMemberStore store.ChatMemberStore
}

type messageContextKey string
const MessageContextKey = messageContextKey("message")

func SetMessageMembership(r *http.Request, message *store.Message) *http.Request{
	ctx := context.WithValue(r.Context(), MessageContextKey, message)
	return r.WithContext(ctx)
}

func GetMessageMembership(r *http.Request) *store.Message {
	messageMember, ok := r.Context().Value(MessageContextKey).(*store.Message)
	if !ok {
		panic("missing message member in request")
	}
	return messageMember
}

func (mm *MessageMiddleware) RequireAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		user, ok := GetUser(r)
		
		if !ok {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error":"signin to continue"})
			return
		}

		messageID, err := utils.ReadParam(r, "msgID")
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid message ID"})
			return
		}

		msg, err := mm.MessageStore.GetMessage(r.Context(), messageID)
		if err != nil {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "message not found"})
			return
		}

		isMember, err := mm.ChatMemberStore.IsMember(r.Context(), msg.ChatID, user.ID)
		if err != nil || !isMember {
			utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error":"access denied"})
			return
		}

		ctx := context.WithValue(r.Context(),"message",msg)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}