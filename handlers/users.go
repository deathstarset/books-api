package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deathstarset/books-api/internal/database"
	"github.com/deathstarset/books-api/mappings"
	"github.com/deathstarset/books-api/responses"
	"github.com/deathstarset/books-api/session"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request, queries *database.Queries) {
	// parse the multi-part data
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, "Failed to parse multiform data")
		return
	}
	// retreiving the file
	file, handler, err := r.FormFile("pfp_link")
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, "Failed to read the file")
		return
	}
	defer file.Close()
	// uploading the file
	fileName, err := uploadFile("uploads/pfp", file, handler)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to upload file into the server : %v", err))
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error : %v", err))
	}

	user, err := queries.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
		Email:     email,
		Username:  username,
		Password:  string(hashedPassword),
		PfpLink:   fileName,
	})
	if err != nil {
		responses.RespondWithError(w, 500, fmt.Sprintf("Create user failed : %v", err))
		return
	}
	responses.RespondWithJSON(w, http.StatusCreated, mappings.DbUserToUserMapping(user))
}

func LoginHandler(client *redis.Client, queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		decoder.Decode(&params)
		user, err := queries.GetUserByUsername(r.Context(), params.Username)
		if err != nil {
			responses.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Failed to get user : %v", err))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
		if err != nil {
			responses.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("Password not match : %v", err))
			return
		}

		redisContext := context.Background()
		sessionObj := session.Session{
			UserID: user.ID,
		}
		sessionJSON, err := json.Marshal(sessionObj)
		if err != nil {
			responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to marshal struct to json : %v", err))
		}
		sessionToken := uuid.NewString()
		err = client.Set(redisContext, sessionToken, sessionJSON, 24*time.Hour*30).Err()
		if err != nil {
			responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to set session cookie : %v", err))
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session_token",
			Value: sessionToken,
		})

		responses.RespondWithJSON(w, http.StatusOK, mappings.DbUserToUserMapping(user))
	}
}

func LogoutHandler(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session_token, err := r.Cookie("session_token")
		if err != nil {
			responses.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Failed to get session cookie : %v", err))
			return
		}
		redisContext := context.Background()
		_, err = client.Del(redisContext, session_token.Value).Result()
		if err != nil {
			responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete session cookie : %v", err))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			MaxAge: -1,
		})
		type responseMessage struct {
			Message string `json:"message"`
		}
		responses.RespondWithJSON(w, http.StatusOK, responseMessage{Message: "Logout successful"})
	}
}
