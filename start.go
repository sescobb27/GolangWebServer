package main

import (
	"log"
	"net/http"
	"webserver"
)

func main() {
	http.HandleFunc("/", webserver.IndexPageHandler)
	http.HandleFunc("/login", webserver.LoginPageHandler)
	http.HandleFunc("/about/", webserver.AboutPageHandler)
	http.HandleFunc("/user/", webserver.UserPageHandler)
	http.HandleFunc("/upload", webserver.UploadPageHandler)
	http.HandleFunc("/signup", webserver.SignUpPageHandler)
	http.HandleFunc("/show", webserver.ShowUploadsHanlder)
	http.HandleFunc("/search", webserver.SearchPageHanlder)
	// To serve a directory on disk (pwd/images) under an alternate URL
	// path (/images/), use StripPrefix to modify the request
	// URL's path before the FileServer sees it:
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
