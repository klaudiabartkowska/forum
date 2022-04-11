package main

import (
	"fmt"
	"net/http"
	"text/template"
	
)

//The "db" package level variable will hold the reference to our database instanc
//var db *sql.DB

var tpl *template.Template // create a container that's  points to the template adress

type User struct {
	username string
	password string
	email    string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func getLoginPage(w http.ResponseWriter, r *http.Request) {

	fmt.Println("login page running")

	tpl.ExecuteTemplate(w, "login.html", nil)

}

func getSignUpPage(w http.ResponseWriter, r *http.Request) {

	tpl.ExecuteTemplate(w, "sign-up.html", nil)

}

func signUpUser(w http.ResponseWriter, r *http.Request) {

}

func loginUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("login user is running now")

	r.ParseForm()

	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println("username:", username)
	fmt.Println("password:", password)

	tpl.ExecuteTemplate(w, "login.html", nil)

}

func homePage(w http.ResponseWriter, r *http.Request) {

	tpl.ExecuteTemplate(w, "homepage.html", nil)
}

func getUser(w http.ResponseWriter, r *http.Request) {

	var u User

	u.username = r.FormValue("username")
	u.password = r.FormValue("password")
	u.email = r.FormValue("email")

	// return User {
	// 	email: u.email,
	// 	password: u.password,
	// }

	tpl.ExecuteTemplate(w, "login.html", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		getLoginPage(w, r)
	case "/sign-up.html":
		getSignUpPage(w, r)
	case "/login-form":
		loginUser(w, r)
	case "/sign-up-form":
		signUpUser(w, r)
	case "/homepage.html":
		homePage(w, r)

	}
}

func main() {

	http.HandleFunc("/", handler)
	fmt.Println("Starting the server on :8080...")

	http.ListenAndServe(":8080", nil)
}
