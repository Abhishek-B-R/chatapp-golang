package routes

import (
	"github.com/Abhishek-B-R/chat-app-golang/internals/app"
	"github.com/go-chi/chi"
)

func SetupRoutes(app *app.Application) *chi.Mux{
	r := chi.NewRouter()

	r.Get("/health",app.HealthCheck)
	r.Post("/auth/login",app.TokenHandler.HandleCreateToken)
	r.Post("/auth/register",app.UserHandler.HandleCreateUser)

	r.Group(func (r chi.Router){
		r.Use(app.UserMiddleware.Authenticate)

		r.Put("/auth/password-reset",app.UserHandler.HandleUpdateUserPassword)
		r.Route("/users", func (r chi.Router){
			r.Get("/me", app.UserHandler.HandleGetCurrentUser)
			r.Put("/me", app.UserHandler.HandleUpdateUser)
			r.Put("/me/last-seen", app.UserHandler.HandleUpdateLastSeen)

			r.Get("/search/{username}", app.UserHandler.HandleGetUserByUsername)
			r.Get("/{userID}", app.UserHandler.HandleGetUserByID)
		})

		r.Route("/chats", func(r chi.Router) {
			r.Get("/", app.ChatHandler.HandleGetUserChats)
			r.Post("/", app.ChatHandler.HandleCreateChat)

			r.Route("/{chatID}", func(r chi.Router) {
				// Middleware: Verify user is member of this chat
				r.Use(app.ChatMiddleware.RequireMembership)

				// Chat details
				r.Get("/", app.ChatHandler.HandleGetChatByID)
				r.Put("/", app.ChatHandler.HandleUpdateChat)
				r.Delete("/", app.ChatHandler.HandleDeleteChat)

				// Chat members management
				r.Route("/members", func(r chi.Router) {
					r.Get("/", app.ChatMemberHandler.HandleGetChatMembers)
					r.Post("/", app.ChatMemberHandler.HandleAddMember)
					r.Delete("/{userID}", app.ChatMemberHandler.HandleRemoveMember)
					r.Get("/{userID}/role", app.ChatMemberHandler.HandleGetUserRole)
					r.Put("/update", app.ChatMemberHandler.HandleUpdateLastRead)
				})

				// Messages in this chats
				r.Route("/messages", func(r chi.Router) {
					r.Get("/{offset}/{limit}", app.MessageHandler.HandleGetChatMessages)
					r.Post("/", app.MessageHandler.HandleCreateMessage)
					r.Get("/unread", app.MessageHandler.HandleGetUnreadCount) 
				})

				// Chat member actions (for current user)
				r.Put("/read", app.ChatMemberHandler.HandleUpdateLastRead)
				r.Put("/mute", app.ChatMemberHandler.HandleMuteChat)
				r.Put("/unmute", app.ChatMemberHandler.HandleUnMuteChat)
			})
		})

		// MESSAGE ROUTES (Individual message operations)
		r.Route("/messages/{msgID}", func(r chi.Router) {
			r.Use(app.MessageMiddleware.RequireAccess)

			r.Get("/", app.MessageHandler.HandleGetMessage)
			r.Put("/", app.MessageHandler.HandleUpdateMessage)
			r.Delete("/", app.MessageHandler.HandleDeleteMessage)
		})




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
		
		r.Post("/message", app.MessageHandler.HandleCreateMessage)
		r.Get("/message/{msgID}", app.MessageHandler.HandleGetMessage)
		r.Get("/messages/{chatID}/{limit}/{offset}", app.MessageHandler.HandleGetChatMessages)
		r.Put("/message/update", app.MessageHandler.HandleUpdateMessage)
		r.Delete("/message/{msgID}", app.MessageHandler.HandleDeleteMessage)
		r.Get("/messages/unread/{chatID}/{userID}",app.MessageHandler.HandleGetUnreadCount)
	})
	return r
}