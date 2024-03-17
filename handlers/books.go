package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/deathstarset/books-api/internal/database"
	"github.com/deathstarset/books-api/mappings"
	"github.com/deathstarset/books-api/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) CreateBookHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, "Failed to parse multiform data")
		return
	}

	file, handler, err := r.FormFile("image_link")
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, "Failed to read the file")
		return
	}
	defer file.Close()

	fileName, err := uploadFile("uploads/books_images", file, handler)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to upload file into the server : %v", err))
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	book, err := apiCfg.Queries.CreateBook(r.Context(), database.CreateBookParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().Local(),
		UpdatedAt:   time.Now().Local(),
		Title:       title,
		Description: description,
		ImageLink:   fileName,
		UserID:      user.ID,
	})
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create book : %v", err))
		return
	}

	responses.RespondWithJSON(w, http.StatusCreated, mappings.DbBookToBookMapping(book))
}

func (apiCfg *ApiConfig) DeleteBookHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	bookIDstr := chi.URLParam(r, "bookID")
	bookID, err := uuid.Parse(bookIDstr)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to parse book id : %v", err))
		return
	}
	err = apiCfg.Queries.DeleteBook(r.Context(), database.DeleteBookParams{
		ID:     bookID,
		UserID: user.ID,
	})
	if err != nil {
		responses.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Failed to delete book : %v", err))
		return
	}
	type deleteResponse struct {
		Message string `json:"message"`
	}
	responses.RespondWithJSON(w, http.StatusAccepted, deleteResponse{Message: "book deleted succefully"})
}

func (apiCfg *ApiConfig) GetAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	DbBooks, err := apiCfg.Queries.GetAllBooks(r.Context())
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get books : %v", err))
	}

	books := []mappings.Book{}
	for _, dbBook := range DbBooks {
		books = append(books, mappings.DbBookToBookMapping(dbBook))
	}

	responses.RespondWithJSON(w, http.StatusOK, books)
}
