package database

import (
	//"database/sql"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"unicode"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template // create a container that's  points to the template adress

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

func (data *Forum) signUpUser(w http.ResponseWriter, r *http.Request) {

	/*  1. chceck e-mail criteria
	    2. check u.username criteria
		 3. check password criteria
		 4. check if u.username is already exists in database
		 5. create bcrypt hash from password
		 6. insert u.username and password hash in database
	*/

	fmt.Println("****Sign-up new user is running ")

	r.ParseForm() // parse the sign-up form

	/************************************ EMAIL ************************************/

	var user User

	user.Email = r.FormValue("email") // grab the email

	var isValidEmail = true

	if isValidEmail != strings.Contains(user.Email, "@") || isValidEmail != strings.Contains(user.Email, ".") {

		isValidEmail = false
	}

	/********************************** USERNAME ******************************/

	user.Username = r.FormValue("username") // grab the u.username (it's a string/ slice of bytes )

	// check user for only alphaNumeric characters
	var isAlphaNumeric = true

	for _, char := range user.Username {
		// func IsLetter(r rune) bool, func IsNumber(r rune) bool
		// if !unicode.IsLetter(char) && if !unicode.IsNumber {              // checking if the char in username are letters and numbers
		if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
			isAlphaNumeric = false

		}

	}

	var nameLength bool

	//checks if name length meets criteria

	if 5 <= len(user.Username) && len(user.Username) <= 50 {
		nameLength = true
	}

	/***************************************** PASSWORD ***********************************/

	// check password criteria
	user.Password = r.FormValue("password")
	//fmt.Println("password:", password, "\n pswdLenght:", len(password))

	// variables that must pass for password creation criteria
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdNoSpaces, pswdLength bool
	pswdNoSpaces = true

	for _, char := range user.Password {
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
	minLenght := 8
	maxLenght := 30

	if minLenght < len(user.Password) && len(user.Password) < maxLenght {
		pswdLength = true
	}
	//fmt.Println("pswdLowercase:", pswdLowercase, "\npswdUppercase:", pswdUppercase, "\npswdNumber:", pswdNumber, "\npswdSpecial:", pswdSpecial, "\npswdLength:", pswdLength, "\npswdNoSpaces:", pswdNoSpaces, "\nnameAlphaNumeric:", isAlphaNumeric, "\nnameLength:", nameLength)
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces || !isAlphaNumeric || !nameLength {
		tpl.ExecuteTemplate(w, "sign-up.html", "please check username and password criteria")
		return
	}

	// check if username already exists for availability

	row := data.DB.QueryRow("SELECT uuid FROM people where username = ?", user.Username)
	var username string
	err := row.Scan(&username)
	if err != sql.ErrNoRows {
		fmt.Println("username already exists. err", err)
		tpl.ExecuteTemplate(w, "sign-up.html", "username already taken")
		return
	}

	row = data.DB.QueryRow("SELECT uuid FROM people where email =?", user.Email)
	var userEmail string
	err = row.Scan(&userEmail)
	if err != sql.ErrNoRows {
		fmt.Println("email is already taken. err", err)
		tpl.ExecuteTemplate(w, "sign-up.html", "email already taken")
	}

	// create hash from password

	var passwordHash []byte

	//func GenerateFromPassword(password []byte, cost int)([]byte, error) generate a password hash from an raw user password,

	passwordHash, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt err:", err)
		tpl.ExecuteTemplate(w, "sign-up.html", "there is a problem registering account")
		return
	}

	sessionID := uuid.NewV4()
	
	data.CreateUser(User{
		Uuid: sessionID.String(),
		Username: user.Username,
		Email:    user.Email,
		Password: string(passwordHash),
	})

	if err != nil {
		tpl.ExecuteTemplate(w, "sign-up.html", "there was a problem registering account")
		return

	// 	//defer .Close()

	} else {
		http.Redirect(w, r, "/login", 302)
	}
	// fmt.Println(user.Username)
	// fmt.Println(user.Email)
	// fmt.Println(passwordHash)

}

func (data *Forum) loginUser(w http.ResponseWriter, r *http.Request) {

 var cookie *http.Cookie

	fmt.Println("*****loginUser is running********")
	
	var user User
	
	r.ParseForm()

	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	
	fmt.Println("username:", user.Username)
	fmt.Println("password:", user.Password)
	
	// retrieve password from db to compare (hash) with user supplied password's hash
	var passwordHash string
	
	row := DB.QueryRow("SELECT password FROM people WHERE Username = ?", user.Username)
	err := row.Scan(&passwordHash)
	fmt.Println("hash from db:", passwordHash)
	if err != nil {
		fmt.Println("error selecting Hash in db by Username")
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}

	// func CompareHashAndPassword(hashedPassword, password []byte) error
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
	// returns nill on succcess
	if err == nil {
		//tpl.ExecuteTemplate(w, "homepage.html",  nil)
		http.Redirect(w, r, "/homepage.html", 302)
		return
	}
	
	fmt.Println("incorrect password")
	tpl.ExecuteTemplate(w, "login.html", "check username and password")
	
		cookie, err = r.Cookie("session")
		if err != nil {
			id := uuid.NewV4()
			//fmt.Println("cookie was not found")
			cookie = &http.Cookie{
				Name:     "session",
				Value:    id.String(),
				Secure:   true,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		}
	
	
}

func (data *Forum) homePage(w http.ResponseWriter, r *http.Request) {




	tpl.ExecuteTemplate(w, "homepage.html", nil)
}

func guestView(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "guest.html", nil)
}

func logout(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Private navbar running ")
	tpl.ExecuteTemplate(writer, "private.navbar.html", nil)
}

func (data *Forum) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		data.homePage(w, r)
	case "/homepage.html":
		data.homePage(w, r)
	case "/login":
		getLoginPage(w, r)
	case "/sign-up.html":
		getSignUpPage(w, r)
	case "/login-form":
		data.loginUser(w, r)
	case "/sign-up-form":
		data.signUpUser(w, r)
	case "/private.navbar.html":
		logout(w, r)
	}
}
