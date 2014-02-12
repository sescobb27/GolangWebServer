package webserver

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"webserver/controllers"
	"webserver/dbconnection"
	"webserver/models"
)

var sessionManager *controllers.SessionManager

func init() {
	mem_pder := controllers.NewMemProvider()
	controllers.RegisterProvider("memory", mem_pder)
	var err error
	sessionManager, err = controllers.InitializeSessionManager("memory", "cookieWeb")
	if err != nil {
		panic(err)
	}
	go sessionManager.GC()
}

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		w.Header().Set("Content-Type", "text/html")
		t, _ := template.ParseFiles("template/index.gtpl")
		t.Execute(w, nil)
	}
}

func AboutPageHandler(w http.ResponseWriter, r *http.Request) {
}

func ShowUploadsHanlder(w http.ResponseWriter, r *http.Request) {
	// session := sessionManager.SessionStart(w, r)
	if r.Method == "GET" {
		r.ParseForm()
		file_arr, _ := dbconnection.GetUsersFiles(10, 0, r.Form.Get("tag"))
		w.Header().Set("Content-Type", "text/html")
		t, _ := template.ParseFiles("template/showfiles.gtpl")
		t.Execute(w, file_arr)
	}
}

func SearchPageHanlder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html")
		tags := dbconnection.GetTags()
		t, _ := template.ParseFiles("template/search.gtpl")
		t.Execute(w, tags)
	}
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	session := sessionManager.SessionStart(w, r)
	if session.Get("username") != "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	switch r.Method {
	case "GET":
		t, _ := template.ParseFiles("template/login.gtpl")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, session.Get("username"))
	case "POST":
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		user := &models.User{Username: &username,
			Password: &password,
		}
		id, err := dbconnection.VerifyUser(user)
		if err != nil {
			fmt.Println(err)
			return
		}
		session.Set("id", id)
		session.Set("username", username)
		http.Redirect(w, r, "/upload", http.StatusFound)
	default:
		fmt.Println("Error on Method: ", r.Method)
	}
}

func UserPageHandler(w http.ResponseWriter, r *http.Request) {
}

func UploadPageHandler(w http.ResponseWriter, r *http.Request) {
	session := sessionManager.SessionStart(w, r)
	if session.Get("username") == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	switch r.Method {
	case "GET":
		crutime := time.Now().Unix()
		hash := md5.New()
		io.WriteString(hash, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", hash.Sum(nil))
		tpl, _ := template.ParseFiles("template/upload.gtpl")
		tpl.Execute(w, token)
	case "POST":
		// 32 << 20 => 33554432 bytes => 32Mb MaxMemory
		r.ParseMultipartForm(32 << 20)
		formfile, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer formfile.Close()

		title := r.Form.Get("title")
		tags := r.Form.Get("tags")

		var file *os.File
		pwd, _ := os.Getwd()
		file_name := pwd + "/images/" + handler.Filename
		file, err = os.OpenFile(file_name, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()
		io.Copy(file, formfile)

		temp := strings.Split(tags, ",")
		tags_arr := make([]string, 0, len(temp))
		for _, tag := range temp {
			tags_arr = append(tags_arr, strings.Trim(tag, " "))
		}
		info, _ := file.Stat()
		user_file := &models.UserFile{Title: title,
			Path:   "/images/" + handler.Filename,
			UserId: session.Get("id").(int64),
			Size:   info.Size(),
			Tags:   tags_arr}
		go dbconnection.InsertUserFile(user_file)
		http.Redirect(w, r, "/search", http.StatusFound)
		return
	default:
		fmt.Println("Error on Method: ", r.Method)
	}
}

func SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	session := sessionManager.SessionStart(w, r)
	if session.Get("username") != "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	switch r.Method {
	case "GET":
		t, _ := template.ParseFiles("template/signup.gtpl")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, "")
	case "POST":
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		user := new(models.User)
		match, err := models.ValidateUsername(username)
		if err != nil {
			fmt.Println(err)
			return
		} else if !match {
			fmt.Println("Validation Error")
			return
		}
		user.Username = &username
		session.Set("username", username)

		match, err = models.ValidatePassword(password)
		if err != nil {
			fmt.Println(err)
			return
		} else if !match {
			fmt.Println("Password Validation Error")
			return
		}
		user.Password = &password

		callback := func(uid int64) {
			session.Set("id", uid)
			fmt.Println(uid)
		}
		err = dbconnection.InsertUser(user, callback)
		if err != nil {
			fmt.Println(err)
			return
		}
		http.Redirect(w, r, "/upload", http.StatusFound)
		return
	default:
		fmt.Fprintf(w, "Action Not Found")
	}
}
