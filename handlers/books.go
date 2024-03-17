package handlers

import (
	"net/http"

	"github.com/deathstarset/books-api/internal/database"
	"github.com/deathstarset/books-api/mappings"
	"github.com/deathstarset/books-api/responses"
)

func (apiCfg *ApiConfig) CreateBookHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	responses.RespondWithJSON(w, http.StatusAccepted, mappings.DbUserToUserMapping(user))
}
