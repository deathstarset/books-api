package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deathstarset/books-api/internal/database"
	"github.com/deathstarset/books-api/mappings"
	"github.com/deathstarset/books-api/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreateBookHandler(w http.ResponseWriter, r *http.Request, user database.User, queries *database.Queries) {

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to read multipart form data : %v", err))
		return
	}

	file, handler, err := r.FormFile("image_link")
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to read file : %v", err))
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
	book, err := queries.CreateBook(r.Context(), database.CreateBookParams{
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

func DeleteBookHandler(w http.ResponseWriter, r *http.Request, user database.User, queries *database.Queries) {
	bookIDstr := chi.URLParam(r, "bookID")
	bookID, err := uuid.Parse(bookIDstr)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to parse book id : %v", err))
		return
	}
	err = queries.DeleteBook(r.Context(), database.DeleteBookParams{
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

func GetAllBooksHandler(w http.ResponseWriter, r *http.Request, queries *database.Queries) {
	DbBooks, err := queries.GetAllBooks(r.Context())
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get books : %v", err))
	}

	books := []mappings.Book{}
	for _, dbBook := range DbBooks {
		books = append(books, mappings.DbBookToBookMapping(dbBook))
	}

	responses.RespondWithJSON(w, http.StatusOK, books)
}

func UpdateBookHandler(w http.ResponseWriter, r *http.Request, user database.User, queries *database.Queries) {
	bookIDstr := chi.URLParam(r, "bookID")
	bookID, err := uuid.Parse(bookIDstr)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to parse uuid : %v", err))
		return
	}
	book, err := queries.GetBookByUserAndID(r.Context(), database.GetBookByUserAndIDParams{
		ID:     bookID,
		UserID: user.ID,
	})
	if err != nil {
		responses.RespondWithJSON(w, http.StatusNotFound, fmt.Sprintf("Book not found : %v", err))
		return
	}
	type parameters struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)
	if params.Description == "" {
		params.Description = book.Description
	}
	if params.Title == "" {
		params.Title = book.Title
	}
	err = queries.UpdateBook(r.Context(), database.UpdateBookParams{
		Title:       params.Title,
		Description: params.Description,
		UpdatedAt:   time.Now().Local(),
		ID:          bookID,
		UserID:      user.ID,
	})
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update book : %v", err))
		return
	}

	type updateBookResponse struct {
		Message string `json:"message"`
	}
	responses.RespondWithJSON(w, http.StatusAccepted, updateBookResponse{Message: "Book updated succefully"})
}
