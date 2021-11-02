package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/appened/note"
	"github.com/gorilla/mux"
)

// TODO: Add logging
// TODO: Add authentication
// TODO: Add marking done
// TODO: Add editing note
// TODO: Add surfacing a note

func main() {
	fmt.Println("Loading folios")
	folios, err := note.LoadFolios()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Loaded %d folios\n", len(folios))

	r := mux.NewRouter()

	// Set up routes
	initailizeRoutes(r, folios)

	// Add middleware
	initailizeMiddleware(r)

	// Start Server
	fmt.Println("Listening on 8081")
	http.ListenAndServe(":8081", r)
}

// Intialize routes
func initailizeRoutes(router *mux.Router, folios []*note.Folio) {
	// GET folios/{name}: Get a folio's notes in an array of strings
	router.HandleFunc("/folios/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		fmt.Printf("GET folios/%s\n", name)

		folio := findFolio(name, folios)
		if folio == nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Folio not found")
			return
		}

		var notes []string
		for _, note := range folio.Notes {
			notes = append(notes, note.Text)
		}
		jsonResponse, err := json.Marshal(notes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}).Methods("GET")

	// POST folios/{name}: Append a note to a folio
	router.HandleFunc("/folios/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		fmt.Printf("POST folios/%s\n", name)

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		folio := findFolio(name, folios)
		if folio == nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Folio not found")
			return
		}

		err := folio.Append(r.FormValue("note"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		fmt.Printf("Created note in folio %v\n", name)
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	// POST folios/: Create a folio
	router.HandleFunc("/folios", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("POST folios/")
		// Get folio name from request
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		name := r.FormValue("name")

		// Validation
		matched, err := regexp.MatchString(`^[a-zA-Z]+$`, name)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !matched {
			fmt.Printf("Bad folio name: %s\n", name)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid folio name, must be one word")
			return
		}
		for _, f := range folios {
			if f.Name == name {
				fmt.Printf("Duplicate folio name: %s\n", name)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Folio with name exists, try a different name")
				return
			}
		}

		// Create New Folio
		folio, err := note.CreateFolio(name)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		folios = append(folios, folio)
		fmt.Printf("Created folio named %v\n", name)

		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")
}

// Initializes Application Middleware
func initailizeMiddleware(router *mux.Router) {
	// Authentication middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get app token
			token := os.Getenv("APPENDED_AUTH_TOKEN")

			// Get client token
			reqToken := r.Header.Get("Authorization")
			splitToken := strings.Split(reqToken, "Bearer ")

			// Check if token was provided and if it is valid
			authorized := len(splitToken) == 2 && splitToken[1] == token

			// Authenticate user
			if authorized {
				next.ServeHTTP(w, r)
			} else {
				// Reject unauthorized requests
				fmt.Println("Unauthorized")
				w.WriteHeader(http.StatusUnauthorized)
			}
		})
	})
}

// Helper for getting a folio. Runs in O(n) time, converting to a map would obviate need for this
func findFolio(name string, folios []*note.Folio) *note.Folio {
	for _, f := range folios {
		if name == f.Name {
			return f
		}
	}
	return nil
}
