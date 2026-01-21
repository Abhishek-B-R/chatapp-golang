package routes

import (
	"github.com/Abhishek-B-R/chat-app-golang/internals/app"
	"github.com/go-chi/chi"
)

func SetupRoutes(app *app.Application) *chi.Mux{
	r := chi.NewRouter()

	r.Group(func (r chi.Router){
		r.Use(app.UserMiddleware.Authenticate)

		r.Post("/chat",app.ChatHandler.HandleCreateChat)
		r.Get("/chat/{chatID}",app.ChatHandler.HandleGetChatByID)
		r.Put("/chat",app.ChatHandler.HandleUpdateChat)
		r.Delete("/chat/{chatID}",app.ChatHandler.HandleDeleteChat)
		r.Get("/user/{userID}/chats",app.ChatHandler.HandleGetUserChats)
		
		r.Get("/user/{userID}", app.UserHandler.HandleGetUserByID)
		r.Get("/user/email/{email}", app.UserHandler.HandlerGetUserByEmail)
		r.Get("/user/username/{username}", app.UserHandler.HandleGetUserByUsername)
		r.Put("/user/update_last_seen/{userID}", app.UserHandler.HandleUpdateLastSeen)
		r.Put("/user",app.UserHandler.HandleUpdateUser)
		
		r.Post("/chat/add", app.ChatMemberHandler.HandleAddMember)
		r.Delete("/chat/delete/{chatID}/{userID}", app.ChatMemberHandler.HandleRemoveMember)
		r.Get("/chat/members/{chatID}", app.ChatMemberHandler.HandleGetChatMembers)
		r.Get("/chat/member/role/{chatID}/{userID}",app.ChatMemberHandler.HandleGetUserRole)
		r.Get("/chat/member/check/{chatID}/{userID}",app.ChatMemberHandler.HandleIsMember)
		r.Put("/chat/update", app.ChatMemberHandler.HandleUpdateLastRead)
		
		r.Post("/message", app.MessageHandler.HandleCreateMessage)
		r.Get("/message/{msgID}", app.MessageHandler.HandleGetMessage)
		r.Get("/messages/{chatID}/{limit}/{offset}", app.MessageHandler.HandleGetChatMessages)
		r.Put("/message/update", app.MessageHandler.HandleUpdateMessage)
		r.Delete("/message/{msgID}", app.MessageHandler.HandleDeleteMessage)
		r.Get("/messages/unread/{chatID}/{userID}",app.MessageHandler.HandleGetUnreadCount)
	})
	
	r.Put("/user/password/{userID}",app.UserHandler.HandleUpdateUserPassword)
	r.Post("/user",app.UserHandler.HandleCreateUser)
	r.Post("/tokens/authentication",app.TokenHandler.HandleCreateToken)
	r.Get("/health",app.HealthCheck)
	return r
}