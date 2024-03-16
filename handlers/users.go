package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/deathstarset/books-api/internal/database"
	"github.com/deathstarset/books-api/mappings"
	"github.com/deathstarset/books-api/responses"
	"github.com/deathstarset/books-api/utils"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
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
	fileName, err := uploadFile("uploads", file, handler)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to upload file into the server : %v", err))
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := apiCfg.Queries.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
		Email:     email,
		Username:  username,
		Password:  password,
		PfpLink:   fileName,
	})
	if err != nil {
		responses.RespondWithError(w, 500, fmt.Sprintf("Create user failed : %v", err))
		return
	}
	responses.RespondWithJSON(w, 201, mappings.DbUserToUserMapping(user))
}

func uploadFile(destination string, file multipart.File, handler *multipart.FileHeader) (string, error) {
	// creating the file destination if not exist
	err := os.MkdirAll(destination, 0755)
	if err != nil {
		return "", err
	}
	// gen random string for file name
	randStr, err := utils.RandomString(10)
	if err != nil {
		return "", err
	}
	// parsing the file extenstion
	extension := strings.Split(handler.Filename, ".")[1]
	// creating a detination path for the file
	destinationPath := filepath.Join(destination, fmt.Sprintf("pfp-%v.%v", randStr, extension))
	// creating the file in ther server
	serverFile, err := os.Create(destinationPath)
	if err != nil {
		return "", err
	}
	defer serverFile.Close()
	// copying the uploaded file into the server file
	if _, err := io.Copy(serverFile, file); err != nil {
		return "", err
	}
	return destinationPath, nil
}
