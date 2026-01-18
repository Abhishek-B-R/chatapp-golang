package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Abhishek-B-R/chat-app-golang/internals/app"
	"github.com/Abhishek-B-R/chat-app-golang/internals/routes"
)

func main() {
	port := 8080

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: r,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	app.Logger.Printf("We are now running on port %d\n", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}