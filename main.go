package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
	//"golang.org/x/crypto/bcrypt"
)

//The "db" package level variable will hold the reference to our database instanc
var db *sql.DB

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

	/*  1. chceck e-mail criteria
	    2. check username criteria
		 3. check password criteria
		 4. check if username is already exists in database
		 5. create bcrypt hash from password
		 6. insert username and password hash in database
	*/

	fmt.Println("****Sign-up new user is running ")

	r.ParseForm() // parse the sign-up form

	/************************************ EMAIL ************************************/

	email := r.FormValue("email") // grab the email

	var isValidEmail = true

	if isValidEmail != strings.Contains(email, "@") {

		isValidEmail = false
	}

	// fmt.Println(isValidEmail)

	//fmt.Println("is valid")

	// if isValidEmail == false
	// 	//fmt.Println("invalid Email")
	// 	tpl.ExecuteTemplate(w, "sign-up.html", nil)
	// 	return
	// }
	/********************************** USERNAME ******************************/

	username := r.FormValue("username") // grab the username (it's a string/ slice of bytes )

	// check user for only alphaNumeric characters
	var isAlphaNumeric = true

	for _, char := range username {
		// func IsLetter(r rune) bool, func IsNumber(r rune) bool
		// if !unicode.IsLetter(char) && if !unicode.IsNumber {              // checking if the char in username are letters and numbers
		if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
			isAlphaNumeric = false

		}

	}

	//	fmt.Println(isAlphaNumeric)

	var nameLength bool

	if 5 <= len(username) && len(username) <= 50 {
		nameLength = true
	}

	//fmt.Println(nameLenght)

	/***************************************** PASSWORD ***********************************/

	// check password criteria
	password := r.FormValue("password")
	//fmt.Println("password:", password, "\n pswdLenght:", len(password))

	// variables that must pass for password creation criteria
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdNoSpaces, pswdLength bool
	pswdNoSpaces = true

	for _, char := range password {
		switch {
		// func IsLower(r rune)bool
		case unicode.IsLower(char):
			pswdLowercase = true
			// func IsUpper(r rune)bool
		case unicode.IsUpper(char):
			pswdUppercase = true
			// func IsNumber(r rune)bool
		case unicode.IsNumber(char):
			pswdNumber = true
			// func IsPunct(r rune)bool, func IsSymbol(r rune)bool
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
			// func IsSpace(r rune) bool, type rune = int32
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false

		}

	}

	if 11 < len(password) && len(password) < 60 {
		pswdLength = true
	}
	//fmt.Println("pswdLowercase:", pswdLowercase, "\npswdUppercase:", pswdUppercase, "\npswdNumber:", pswdNumber, "\npswdSpecial:", pswdSpecial, "\npswdLength:", pswdLength, "\npswdNoSpaces:", pswdNoSpaces, "\nnameAlphaNumeric:", isAlphaNumeric, "\nnameLength:", nameLength)
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces || !isAlphaNumeric || !nameLength {
		tpl.ExecuteTemplate(w, "sign-up.html", "wrong")
		return
	}

	// db, err := sql.Open("sqlite3", "userdata.db")
	// if err != nil {
	// 	fmt.Println(err.Error())

	// }
		fmt.Fprint(w, "congrats, your account has been successfully created")
	}


// fmt.Println("email:", email)
// fmt.Println("username:", username)
// fmt.Println("password:", password)

func loginUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("*****loginuser is running********")

	r.ParseForm()

	username := r.FormValue("username")

	if len(username) >= 2 && len(username) <= 8 {
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
	case "/goback.html":
		signUpUser(w, r)
	}
}

func main() {

	var err error
	db, err = sql.Open("sqlite3", "userdata.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/", handler)
	fmt.Println("Starting the server on :8080...")
	http.ListenAndServe(":8080", nil)

}
