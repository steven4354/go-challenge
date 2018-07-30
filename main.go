package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"html/template"
	"bytes"
)

// controller for / route 
func initialController(w http.ResponseWriter, r *http.Request) {
	// create a struct that holds information to be displayed in our HTML file
	type TemplateData struct {}
	
	templates := template.Must(template.ParseFiles("templates/welcome-template.html"))

	templateData := TemplateData{}
	
	// serves the first view file and checks for errors
	 if err := templates.ExecuteTemplate(w, "welcome-template.html", templateData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
	 }
}

// controller for balances route
func balancesController(w http.ResponseWriter, r *http.Request) {

	// @section1: reading the response body
	type Balance struct{
		Address string
		Contract string
	}

	var b Balance

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&b)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	
	fmt.Println("Post request recieved. It's payload ->")
	fmt.Println("Address:" + b.Address);
	fmt.Println("Contract:" + b.Contract);

	// @section2: creating a post request to our heroku api
	x := new(bytes.Buffer)
	json.NewEncoder(x).Encode(b)
	res, _ := http.Post("https://web3-challenge-heroku.herokuapp.com/balances", "application/json; charset=utf-8", x)
	
	type BalancesAPIResponse struct{
		Balance string
	}

	var bar BalancesAPIResponse

	json.NewDecoder(res.Body).Decode(&bar)

	fmt.Println("Post request to heroku working. It's response payload ->")
	fmt.Println("Balance:" + bar.Balance)

	// @section3: sending back to the user a JSON payload
	json.NewEncoder(w).Encode(bar)
}

func main() {	
	// serve the static files
	http.Handle("/static/",
		 http.StripPrefix("/static/",
				http.FileServer(http.Dir("static")))) 

	// initial route serves the first view file
	http.HandleFunc("/" , initialController)
	http.HandleFunc("/balances" , balancesController)

	fmt.Println("Listening");
	fmt.Println(http.ListenAndServe(":8080", nil));
}