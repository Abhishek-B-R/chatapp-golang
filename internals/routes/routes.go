package routes

import (
	"github.com/Abhishek-B-R/chat-app-golang/internals/app"
	"github.com/go-chi/chi"
)

func SetupRoutes(app *app.Application) *chi.Mux{
	r := chi.NewRouter()

	r.Post("/chat",app.ChatHandler.HandleCreateChat)
	r.Get("/chat/{chatID}",app.ChatHandler.HandleGetChatByID)
	r.Put("/chat",app.ChatHandler.HandleUpdateChat)
	r.Delete("/chat/{chatID}",app.ChatHandler.HandleDeleteChat)
	
	r.Get("/user/{userID}/chats",app.ChatHandler.HandleGetUserChats)
	r.Post("/user",app.UserHandler.HandleCreateUser)

	r.Post("/chat/members/add", app.ChatMemberHandler.HandleAddMember)

	r.Get("/health",app.HealthCheck)
	return r
}