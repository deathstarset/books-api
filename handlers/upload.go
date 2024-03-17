package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/deathstarset/books-api/utils"
)

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
