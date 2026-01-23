package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Abhishek-B-R/chat-app-golang/internals/api"
	"github.com/Abhishek-B-R/chat-app-golang/internals/middleware"
	"github.com/Abhishek-B-R/chat-app-golang/internals/store"
	"github.com/Abhishek-B-R/chat-app-golang/migrations"
)

type Application struct{
	Logger *log.Logger
	ChatHandler *api.ChatHandler
	MessageHandler *api.MessageHandler
	ChatMemberHandler *api.ChatMemberHandler
	UserHandler *api.UserHandler
	TokenHandler *api.TokenHandler
	
	UserMiddleware middleware.UserMiddleware
	ChatMiddleware middleware.ChatMiddleware
	MessageMiddleware middleware.MessageMiddleware
	DB *sql.DB
}

func NewApplication() (*Application, error){
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	chatStore := store.NewPostgresChatStore(pgDB)
	messageStore := store.NewPostgresMessageStore(pgDB)
	chatMemberStore := store.NewPostgresChatMemberStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	chatHandler := api.NewChatHandler(chatStore, logger)
	messageHandler := api.NewMessageHandler(messageStore, logger)
	chatMemberHandler := api.NewChatMemberHandler(chatMemberStore, messageStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	userMiddlewareHandler := middleware.UserMiddleware{UserStore: userStore}
	chatMiddlewareHandler := middleware.ChatMiddleware{ChatMemberStore: chatMemberStore}
	messageMiddlewareHandler := middleware.MessageMiddleware{MessageStore: messageStore, ChatMemberStore: chatMemberStore}

	app := &Application{
		Logger: logger,
		ChatHandler: chatHandler,
		MessageHandler: messageHandler,
		ChatMemberHandler: chatMemberHandler,
		UserHandler: userHandler,
		TokenHandler: tokenHandler,
		UserMiddleware: userMiddlewareHandler,
		ChatMiddleware: chatMiddlewareHandler,
		MessageMiddleware: messageMiddlewareHandler,
		DB: pgDB,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is working pretty fine")
}

