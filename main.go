package main

import (
	"net/http"
	"fmt"
	"html/template"
)

//Create a struct that holds information to be displayed in our HTML file
type TemplateData struct {}

func main() {	
	templates := template.Must(template.ParseFiles("templates/welcome-template.html"))

	templateData := TemplateData{}

	// serve the static files
	http.Handle("/static/",
		 http.StripPrefix("/static/",
				http.FileServer(http.Dir("static")))) 

	// initial route serves the first view file
	http.HandleFunc("/" , func(w http.ResponseWriter, r *http.Request) {
		// serves the first view file and checks for errors
		 if err := templates.ExecuteTemplate(w, "welcome-template.html", templateData); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
		 }
	})

	fmt.Println("Listening");
	fmt.Println(http.ListenAndServe(":8080", nil));
}