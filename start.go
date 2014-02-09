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
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
