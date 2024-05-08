package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// API endpoints
var (
	artistsEndpoint = "https://groupietrackers.herokuapp.com/api/artists"
)

// Structs to unmarshal API responses
type Artist struct {
	Id           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Year         int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Members      []string `json:"members"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDate"`
	Relations    string   `json:"relations"`
}

type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Error struct {
	ErrorNum int
	ErrorMsg string
}

// Handler function to fetch artists data
func homePage(w http.ResponseWriter, r *http.Request) {
	// if r.Method == "GET" && r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }
	var artistsData []Artist
	// Make GET request to artists API endpoint
	resp, err := http.Get(artistsEndpoint)
	if err != nil {
		http.Error(w, "Failed to fetch artists data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode JSON response
	err = json.NewDecoder(resp.Body).Decode(&artistsData)
	if err != nil {
		http.Error(w, "Failed to decode artists data", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	err = tmpl.Execute(w, artistsData)
	if err != nil {
		http.Error(w, "Failed to execute file index.html", http.StatusInternalServerError)
		return
	}
}

func relationPage(w http.ResponseWriter, r *http.Request) {
	// if r.Method == "GET" && r.URL.Path != "/relation" {
	// 	http.NotFound(w, r)
	// 	return
	// }
	relationsEndpoint := r.FormValue("relationlink")
	fmt.Println(relationsEndpoint)
	var relationData Relation
	resp, err := http.Get(relationsEndpoint)
	if err != nil {
		http.Error(w, "Failed to fetch relation data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode JSON response
	err = json.NewDecoder(resp.Body).Decode(&relationData)
	if err != nil {
		http.Error(w, "Failed to decode artists data", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("template/relation.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	if r.URL.Path != "/relation" {
		ErrorHandler(w, r, http.StatusNotFound)
	}
	err = tmpl.Execute(w, relationData)
	if err != nil {
		http.Error(w, "Failed to execute file relation.html", http.StatusInternalServerError)
		return
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err int) {
	temp, errorCode := template.ParseFiles("template/error.html")
	if errorCode != nil {
		log.Fatal(errorCode)
		return
	}
	w.WriteHeader(err)
	errorData := Error{ErrorNum: err}

	if err == 404 {
		errorData.ErrorMsg = "Page Not Found"

	} else if err == 500 {
		errorData.ErrorMsg = "Internal Server Error"
	} else if err == 400 {
		errorData.ErrorMsg = "Bad Request"
	}
	temp.Execute(w, errorData)

}

func main() {
	styles := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", styles))

	// Define HTTP routes
	http.HandleFunc("/", homePage)
	http.HandleFunc("/relation", relationPage)

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
