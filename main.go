package main

import (
	"database/sql"
	"fmt"

	"net/http"
	"strings"
	"text/template"
	"unicode"
	"github.com/satori/go.uuid"

	"example.com/m/database"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

//The "db" package level variable will hold the reference to our database instanc
///var DB *sql.DB

//var insertStmt *sql.Stmt

var tpl *template.Template // create a container that's  points to the template adress

type User struct {
	Username string
	Password string
	Email    string
}


func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}


// getLoginPage serves form for log in existing users

func getLoginPage(w http.ResponseWriter, r *http.Request) {

	fmt.Println("*****loginpage running****")

	tpl.ExecuteTemplate(w, "login.html", nil)
	// t := parseTemplateFiles("login")
	// t.Execute(w,nil)

}

// getSingUpPage serves form for signing up new users

func getSignUpPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("***sign-up page runnning***")
	tpl.ExecuteTemplate(w, "sign-up.html", nil)

}

func signUpUser(w http.ResponseWriter, r *http.Request) {

	/*  1. chceck e-mail criteria
	    2. check u.username criteria
		 3. check password criteria
		 4. check if u.username is already exists in database
		 5. create bcrypt hash from password
		 6. insert u.username and password hash in database
	*/

	fmt.Println("****Sign-up new user is running ")

	var u User
   

	r.ParseForm() // parse the sign-up form

	/************************************ EMAIL ************************************/

	u.Email = r.FormValue("email") // grab the email

	var isValidEmail = true

	if isValidEmail != strings.Contains(u.Email, "@") {

		isValidEmail = false
	}

	/********************************** USERNAME ******************************/

	u.Username = r.FormValue("username") // grab the u.username (it's a string/ slice of bytes )

	// check user for only alphaNumeric characters
	var isAlphaNumeric = true

	for _, char := range u.Username {
		// func IsLetter(r rune) bool, func IsNumber(r rune) bool
		// if !unicode.IsLetter(char) && if !unicode.IsNumber {              // checking if the char in username are letters and numbers
		if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
			isAlphaNumeric = false

		}

	}

	var nameLength bool

	if 5 <= len(u.Username) && len(u.Username) <= 50 {
		nameLength = true
	}

	/***************************************** PASSWORD ***********************************/

	// check password criteria
	u.Password = r.FormValue("password")
	//fmt.Println("password:", password, "\n pswdLenght:", len(password))

	// variables that must pass for password creation criteria
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdNoSpaces, pswdLength bool
	pswdNoSpaces = true

	for _, char := range u.Password {
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

	if 11 < len(u.Password) && len(u.Password) < 60 {
		pswdLength = true
	}
	//fmt.Println("pswdLowercase:", pswdLowercase, "\npswdUppercase:", pswdUppercase, "\npswdNumber:", pswdNumber, "\npswdSpecial:", pswdSpecial, "\npswdLength:", pswdLength, "\npswdNoSpaces:", pswdNoSpaces, "\nnameAlphaNumeric:", isAlphaNumeric, "\nnameLength:", nameLength)
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces || !isAlphaNumeric || !nameLength {
		tpl.ExecuteTemplate(w, "sign-up.html", "please check username and password criteria")
		return
	}

	stmt := "SELECT id FROM people where username = ?"
	row := database.DB.QueryRow(stmt, u.Username)

	var id string
	err := row.Scan(&id)
	if err != sql.ErrNoRows {
		fmt.Println("username already exists. err", err)
		tpl.ExecuteTemplate(w, "sign-up.html", "username already taken")
		return
	}

	stmt = "SELECT id FROM people where email = ?"
	row = database.DB.QueryRow(stmt, u.Username)

	var userEmail string
	err = row.Scan(&userEmail)
	if err != sql.ErrNoRows {
		fmt.Println("email is already taken. err", err)
		tpl.ExecuteTemplate(w, "sign-up.html", "email already taken")
		return
	}


	// create hash from password
	var passwordHash []byte

	// func GenerateFromPassword(password []byte, cost int)([]byte, error) generate a password hash from an raw user password,

	passwordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt err:", err)
		tpl.ExecuteTemplate(w, "sign-up.html", "there is a problem registering account")
		return
	}

	var insertStmt *sql.Stmt
	insertStmt, err = database.DB.Prepare("INSERT INTO people (username, email, passwordHash) VALUES (?, ?, ?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "sign-up.html", "there was a problem registering account")
		return
	}
	defer insertStmt.Close()

	var result sql.Result
	result, err = insertStmt.Exec(u.Username, u.Email, passwordHash)
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	fmt.Println("err:", err)
	if err != nil {
		fmt.Println("error inserting new user")
		tpl.ExecuteTemplate(w, "sing-up.html", "there was a problem registering account")
		return
	} else {

		fmt.Println("hash:", passwordHash)
		fmt.Println("string(hash)", string(passwordHash))

		//tpl.ExecuteTemplate(w, "homepage.html", "congrats, your account has been successfully created")
		http.Redirect(w, r, "/login",302)
	}
}

// fmt.Println("email:", email)
// fmt.Println("u.username:", u.username)
// fmt.Println("password:", password)

func loginUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("*****loginUser is running********")

	cookie, err := r.Cookie("session")
	if err != nil {
		id := uuid.NewV4()
		//fmt.Println("cookie was not found")
		cookie = &http.Cookie{
			Name: "session",
			Value: id.String(),
			Secure: true, 
			HttpOnly: true, 
		}
		http.SetCookie(w, cookie)
	}
	fmt.Println("cookie:",cookie)

	var u User

	r.ParseForm()

	u.Username = r.FormValue("username")
	u.Password = r.FormValue("password")

	fmt.Println("username:", u.Username)
	fmt.Println("password:", u.Password)

	// retrieve password from db to compare (hash) with user supplied password's hash
	var passwordHash string
	stmt := "SELECT passwordHash FROM people WHERE Username = ?"
	row := database.DB.QueryRow(stmt, u.Username)
	err = row.Scan(&passwordHash)
	fmt.Println("hash from db:", passwordHash)
	if err != nil {
		fmt.Println("error selecting Hash in db by Username")
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}
	// func CompareHashAndPassword(hashedPassword, password []byte) error
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(u.Password))
	// returns nill on succcess
	if err == nil {
		//tpl.ExecuteTemplate(w, "homepage.html",  nil)
		http.Redirect(w, r, "/homepage.html", 302)
		return
	}

	fmt.Println("incorrect password")
	tpl.ExecuteTemplate(w, "login.html", "check username and password")

}


func homePage(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "homepage.html", nil)
}

func guestView(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "guest.html", nil)
}


func logout(writer http.ResponseWriter, request *http.Request) {


	http.Redirect(writer, request, "/", 302)
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
	case "/guest.html":
		guestView(w, r)
	}
}

func main() {
	database.Connect()
	http.HandleFunc("/", handler)
	
	fmt.Println("Starting the server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
