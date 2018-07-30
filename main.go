package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"html/template"
)

// create a struct that holds information to be displayed in our HTML file
type TemplateData struct {}

// struct for json responses
type JsonResponse struct {
	// Reserved field to add some meta information to the API response
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

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

	// test post route
	http.HandleFunc("/test-post" , func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response := &JsonResponse{Data: "hello"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
	})

	fmt.Println("Listening");
	fmt.Println(http.ListenAndServe(":8080", nil));
}