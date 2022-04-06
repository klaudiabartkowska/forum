package main

import (
	"net/http"
	"text/template"
)

var tpl *template.Template // create a container that's  points to the template adress

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}


func getLoginPage(w http.ResponseWriter, r *http.Request){



 tpl.ExecuteTemplate(w, "/login", )

}



func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		getLoginPage(w,r)
	// case "/sign-up.html":
	// 	getSignUpPage(w,r)
	}

}

func main() {

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
