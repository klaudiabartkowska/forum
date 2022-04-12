package main

import (
	//"database/sql"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"unicode"
	//"github.com/mattn/go-sqlite3"
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

// getLoginPage serves form for log in existing users  

func getLoginPage(w http.ResponseWriter, r *http.Request) {

	fmt.Println("*****loginpage running****")

	tpl.ExecuteTemplate(w, "login.html", nil)

}

// getSingUpPage serves form for signing up new users 

func getSignUpPage(w http.ResponseWriter, r *http.Request) {
   fmt.Println("***sign-up page runnning***")
	tpl.ExecuteTemplate(w, "sign-up.html", nil)

}

func signUpUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("****Sing-up new user is running " )

	r.ParseForm()

	email := r.FormValue("email")


	isValidEmail := strings.Contains(email,"@")

	if isValidEmail {
		fmt.Println("is valid")
		
	}else {
		fmt.Println("invalid Emial")
		tpl.ExecuteTemplate(w, "sign-up.html", nil)
		return
	}
	
	username := r.FormValue("username")


// check user for only alphaNumeric characters 

  var nameAlphaNumeric  = true 

  for _, char := range username {
	  if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
		  nameAlphaNumeric = false 
	  }
  
	}
	fmt.Print(nameAlphaNumeric)

	var isValidLenght bool  
	
	if len(username) <= 5 && len(username) >= 8 {
		isValidLenght = true 
	} 

fmt.Println(isValidLenght)
	password := r.FormValue("password")
	
	
	
	fmt.Println("email:", email)
	fmt.Println("username:", username)
	fmt.Println("password:", password)
	


}



func loginUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("*****loginuser is running********")

	r.ParseForm()

	username := r.FormValue("username")

	if len(username) < 8 {
		fmt.Println("Username is too short")
		tpl.ExecuteTemplate(w, "login.html", nil)
		return 
	}


	password := r.FormValue("password")

	fmt.Println("username:", username)
	fmt.Println("password:", password)


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
