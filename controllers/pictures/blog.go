package pictures

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"photo_blog/views"
	"strings"
)

func Blog() http.HandlerFunc {
	return blog
}

func blog(w http.ResponseWriter, req *http.Request) {
	picsPath, err := GetPicsPath()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Blog page:", err)
		return
	}

	userDirs, err := ioutil.ReadDir(picsPath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Blog page, readdir:", err)
		return
	}

	userNames := make([]string, 0, len(userDirs))
	for _, userDir := range userDirs {
		if userDir.IsDir() {
			userNames = append(userNames, userDir.Name())
		}
	}

	if err = views.Tpl().ExecuteTemplate(w, "blogList.gohtml", userNames); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Blog page. ExecuteTemplate:", err)
		return
	}
}

func GetPicsPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	picsPath := filepath.Join(wd, UserPicturesDir)
	return picsPath, nil
}

func UserBlog() http.HandlerFunc {
	return userBlog
}

func userBlog(w http.ResponseWriter, req *http.Request) {
	uriParts := strings.Split(req.RequestURI, "/")
	if len(uriParts) != 3 {
		http.Error(w, "Not found", http.StatusNotFound)
		log.Println("UserBlog page. Trying to get some shit:", req.RequestURI)
		return
	}

	picsPath, err := GetPicsPath()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Blog page:", err)
		return
	}

	userName := uriParts[2]
	userPath := filepath.Join(picsPath, userName)

	// check is it dir and is it exists
	if userDirInfo, err := os.Stat(userPath); err != nil || !userDirInfo.IsDir() {
		http.Error(w, "Not found", http.StatusNotFound)
		log.Println("UserBlog page. Trying to get some shit:", err)
		if userDirInfo != nil {
			log.Println(userDirInfo.Name())
		}
		return
	}

	pics, err := ioutil.ReadDir(userPath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("User blog page. Get pics:", err)
		return
	}

	picsHrefs := make([]string, 0, len(pics))
	for _, pic := range pics {
		picsHrefs = append(picsHrefs, filepath.Join("/pictures", userName, pic.Name()))
	}

	blogInfo := struct {
		UserName  string
		PicsHrefs []string
	}{
		UserName:  userName,
		PicsHrefs: picsHrefs,
	}

	if err = views.Tpl().ExecuteTemplate(w, "blog.gohtml", blogInfo); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("User Blog page. ExecuteTemplate:", err)
		return
	}
}
