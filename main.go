package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// Define the struct to map the JSON response
type Pokemon struct {
	Name           string  `json:"name"`
	BaseExperience int     `json:"base_experience"`
	Height         int     `json:"height"`
	Weight         int     `json:"weight"`
	Sprites        Sprites `json:"sprites"`
	Types          []Type  `json:"types"`
}

type Sprites struct {
	Other        OtherSprites `json:"other"`
	FrontDefault string       `json:"front_default"`
}

type OtherSprites struct {
	Showdown Showdown `json:"showdown"`
}

type Showdown struct {
	FrontDefault string `json:"front_default"`
}

type Type struct {
	Slot       int        `json:"slot"`
	TypeDetail TypeDetail `json:"type"`
}

type TypeDetail struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

//go:embed views/*
var views embed.FS
var t = template.Must(template.ParseFS(views, "views/*"))

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	})

	router.HandleFunc("POST /poke", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusInternalServerError)
		}

		resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(r.FormValue("pokemon")))
		if err != nil {
			http.Error(w, "Unable to fetch new pokemon", http.StatusInternalServerError)
		}
		data := Pokemon{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, "Unable to parse the Pokemon data", http.StatusInternalServerError)
		}
		if err := t.ExecuteTemplate(w, "response.html", data); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	})

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	fmt.Println("Listening on 3000")
	server.ListenAndServe()
}
