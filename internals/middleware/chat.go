package middleware

// import (
// 	"net/http"

// 	"github.com/Abhishek-B-R/chat-app-golang/internals/store"
// 	"github.com/Abhishek-B-R/chat-app-golang/internals/utils"
// )

// type ChatMiddleware struct{
// 	chatMemberStore store.ChatMemberStore
// }

// func NewChatMiddleware(cms store.ChatMemberStore) *ChatMiddleware {
// 	return &ChatMiddleware{chatMemberStore: cms}
// }

// func (m *ChatMiddleware) RequireMembership(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
// 		user := r.Context().Value("user").(*store.User)
// 		chatID, err := utils.ReadParam(r, "chatID")

// 		if err != nil {
// 			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error":"invalid chat ID"})
// 			return
// 		}

// 				// Check if user is a member
// 		isMember, err := m.chatMemberStore.IsMember(chatID, user.ID)
// 		if err != nil || !isMember {
// 			utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you are not a member of this chat"})
// 			return
// 		}
		
// 		// Add chatID to context for convenience
// 		ctx := context.WithValue("chatID", chatID)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }