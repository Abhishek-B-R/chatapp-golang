package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Abhishek-B-R/chat-app-golang/internals/store"
	"github.com/Abhishek-B-R/chat-app-golang/internals/utils"
)

type UserMiddleware struct{
	UserStore store.UserStore
}

type contextKey string
const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) (*store.User, bool) {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	return user, ok && user != nil
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		// in this anonymous fn, we can interject any incoming requests to our server

		w.Header().Add("Vary","Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error":"missing authorization header"})
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error":"invalid authorization header"})
			return
		}
		token := parts[1]

		user, err := um.UserStore.GetUserToken(r.Context(), token)
		if err != nil {
			fmt.Println(err)
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error":"invalid or expired token"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}