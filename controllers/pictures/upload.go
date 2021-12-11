package pictures

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"photo_blog/components/auth"
	"photo_blog/views"
	"strings"
)

const UserPicturesDir = "contentStore"

func Upload() http.HandlerFunc {
	return upload
}

func upload(w http.ResponseWriter, req *http.Request) {
	uploadInfo := struct {
		IsUploaded bool
		FileName   string
	}{req.Method == http.MethodPost, ""}

	if req.Method == http.MethodPost {
		if uploadInfo.FileName = handleUpload(w, req); uploadInfo.FileName == "" {
			return
		}
	}

	if err := views.Tpl().ExecuteTemplate(w, "upload.gohtml", uploadInfo); err != nil {
		log.Println(err)
	}
}

func handleUpload(w http.ResponseWriter, req *http.Request) string {
	file, header, err := req.FormFile("pic")
	if err != nil {
		http.Error(w, "Error while uploading file", http.StatusInternalServerError)
		log.Println("Error while uploading file:", err)
		return ""
	}
	defer file.Close()

	fileName, err := generateFileName(header.Filename, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return ""
	}

	authUser := auth.GetUserFromSession(req)
	if err = createPicFile(fileName, file, authUser.Login); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error createPicFile:", err)
		return ""
	}

	return fileName
}

var ErrBadExtension = errors.New("bad file extension")

func generateFileName(uplFileName string, uplFile multipart.File) (string, error) {
	nameParts := strings.Split(uplFileName, ".")
	if len(nameParts) < 2 {
		return "", ErrBadExtension
	}

	ext := nameParts[1]
	hash := sha1.New()
	if _, err := io.Copy(hash, uplFile); err != nil {
		return "", err
	}
	// return file offset to begin
	if _, err := uplFile.Seek(0, 0); err != nil {
		return "", err
	}

	newName := fmt.Sprintf("%x.%s", hash.Sum(nil), ext)
	return newName, nil
}

func createPicFile(fileName string, uplFile multipart.File, userName string) error {
	userDir, err := createUserDirIfNeeded(userName)
	if err != nil {
		return err
	}

	path := filepath.Join(userDir, fileName)
	newFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer newFile.Close()

	//uplFile.Seek(0, 0)
	if _, err := io.Copy(newFile, uplFile); err != nil {
		return err
	}

	return nil
}

func createUserDirIfNeeded(userName string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	userPath := filepath.Join(wd, UserPicturesDir, userName)
	if err := os.MkdirAll(userPath, os.ModePerm); err != nil {
		return "", err
	}

	return userPath, nil
}
