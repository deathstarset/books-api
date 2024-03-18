package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deathstarset/books-api/internal/database"
	"github.com/deathstarset/books-api/responses"
	"github.com/deathstarset/books-api/session"
	"github.com/redis/go-redis/v9"
)

type authHandler func(w http.ResponseWriter, r *http.Request, user database.User, queries *database.Queries)

func AuthMiddleware(handler authHandler, client *redis.Client, queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the session id from the cookies
		session_cookie, err := r.Cookie("session_token")
		if err != nil {
			responses.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("User not authorized : %v", err))
			return
		}
		session_token := session_cookie.Value
		redisContext := context.Background()
		userInfo, err := client.Get(redisContext, session_token).Result()
		if err != nil {
			responses.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("User not authorized : %v", err))
			return
		}
		session := session.Session{}
		err = json.Unmarshal([]byte(userInfo), &session)
		if err != nil {
			responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to decode string into json : %v", err))
			return
		}
		userId := session.UserID
		user, err := queries.GetUserByID(r.Context(), userId)
		if err != nil {
			responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get user : %v", err))
			return
		}
		handler(w, r, user, queries)
	}
}

type handlerFunc func(w http.ResponseWriter, r *http.Request, queries *database.Queries)

func QueriesMiddleware(handler handlerFunc, queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, queries)
	}
}
