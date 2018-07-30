package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"html/template"
	"bytes"
	"io/ioutil"
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

	// posting to to retrieve balances info https://web3-challenge-heroku.herokuapp.com/balances
	http.HandleFunc("/balances" , func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request info", r)

		var jsonStr = []byte(`{}`)
		req, error := http.NewRequest("POST", "https://web3-challenge-heroku.herokuapp.com/balances", bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, error := client.Do(req)
    if error != nil {
        panic(error)
    }
		defer resp.Body.Close()
		
		fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))

		//orig
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response := &JsonResponse{Data: string(body)}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/testbalances" , func(w http.ResponseWriter, r *http.Request) {

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
		
		// @section2: creating a post request to our heroku api
		var jsonStr = []byte(`{"Address":"` + b.Address + `",` + `"Contract":"` + b.Contract + `"}`)
		req, error := http.NewRequest("POST", "https://web3-challenge-heroku.herokuapp.com/balances", bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, error := client.Do(req)
    if error != nil {
        panic(error)
    }
		defer resp.Body.Close()
		
		fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	})

	fmt.Println("Listening");
	fmt.Println(http.ListenAndServe(":8080", nil));
}