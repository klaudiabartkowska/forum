package main

import (
	"fmt"
	"net/http"
	"text/template"
)

var tpl *template.Template // create a container that's  points to the template adress

type User struct {
	username string
	password string
	email string
}

// func templating(w http.ResponseWriter, fileName string, data interface{}) {

// 	t, _ := template.ParseFiles(fileName)
// 	t.ExecuteTemplate(w, fileName, data)

// }

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}



func getLoginPage(w http.ResponseWriter, r *http.Request) {

	// if r.URL.Path != "/"{
	// 	http.NotFound(w, r)
	// 	return
	// }

	tpl.ExecuteTemplate(w, "login.html", nil)
}


func getSignUpPage(w http.ResponseWriter, r *http.Request) {

  tpl.ExecuteTemplate(w,"sign-up.html",nil )

}


func signUpUser(w http.ResponseWriter, r *http.Request){

}

func loginUser(w http.ResponseWriter, r *http.Request){

}


func homePage(w http.ResponseWriter, r *http.Request){


	tpl.ExecuteTemplate(w,"homepage.html",nil)
}


func getUser(w http.ResponseWriter, r *http.Request){

   var u User

	u.username = r.FormValue("username")
	u.password = r.FormValue("password")
	u.email = r.FormValue("email")


	fmt.Println(u)

	tpl.ExecuteTemplate(w, "login.html",u)

}




func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {	
	case "/login":
		getLoginPage(w, r)
	case "/sign-up.html":
	 	getSignUpPage(w, r)
	case "/login-form":
		loginUser(w,r)
	case "/sign-up-form":
		signUpUser(w,r)
	case "/homepage.html":
		homePage(w,r)		

	}

}

func main() {

	http.HandleFunc("/", handler)
	fmt.Println("Starting the server on :8080...")
	http.ListenAndServe(":8080", nil)
}
