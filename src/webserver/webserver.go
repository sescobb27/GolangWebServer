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
	r.ParseForm()                 // parse arguments, you have to call this by yourself
	fmt.Println("form: ", r.Form) // print form information in server side
	fmt.Println("path: ", r.URL.Path)
	fmt.Println("scheme: ", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	fmt.Println("host: ", r.Host)
	fmt.Println("header: ", r.Header)
	fmt.Println("method: ", r.Method)
	fmt.Println("requestUri: ", r.URL.RequestURI())
	fmt.Println("requestUri: ", r.URL.String())
	fmt.Println("rawQuery: ", r.URL.RawQuery)
	fmt.Println("urlHost: ", r.URL.Host)
	fmt.Println("urlFragment: ", r.URL.Fragment)
	fmt.Println("urlScheme: ", r.URL.Scheme)
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello From Golang")
}

func AboutPageHandler(w http.ResponseWriter, r *http.Request) {
}

func ShowUploadsHanlder(w http.ResponseWriter, r *http.Request) {
	session := sessionManager.SessionStart(w, r)

}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	session := sessionManager.SessionStart(w, r)
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
		session.Set("username", username)
		id, err := dbconnection.VerifyUser(user)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("user id: ", id)
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
		fmt.Fprintf(w, "%v", handler.Header)

		var file *os.File
		file_name := "images/" + handler.Filename
		file, err = os.OpenFile(file_name, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()
		io.Copy(file, formfile)

		info, _ := file.Stat()
		user_file := &models.UserFile{Title: title,
			Path:   file.Name(),
			UserId: session.Get("id").(int64),
			Size:   info.Size()}
		go dbconnection.InsertUserFile(user_file)
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
		t, _ := template.ParseFiles("template/login.gtpl")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, nil)
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
