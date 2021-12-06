package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/appened/HTTPLogger"
	"github.com/appened/note"
	"github.com/gorilla/mux"
)

// TODO: Add editing note
// TODO: Add surfacing a note

func main() {
	// Init Logger
	logger := HTTPLogger.New(os.Stdout, HTTPLogger.LOG_ALL)

	// Load Folios
	logger.Info("Loading folios")
	folios, err := note.LoadFolios()
	if err != nil {
		logger.Error(err)
	}
	logger.Info(fmt.Sprintf("Loaded %d folios\n", len(folios)))

	r := mux.NewRouter()

	// Add middleware
	initailizeMiddleware(r, logger)

	// Set up routes
	initailizeRoutes(r, logger, folios)

	// Start Server
	logger.Info("Listening on 8081")
	http.ListenAndServe(":8081", r)
}

// Intialize routes
func initailizeRoutes(router *mux.Router, logger *HTTPLogger.Logger, folios map[string]*note.Folio) {
	// GET folios/{name}: Get a folio's notes in an array of strings
	router.HandleFunc("/folios/{name}{slash:/?}", func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]

		folio := folios[name]
		if folio == nil {
			w.WriteHeader(http.StatusNotFound)
			logger.InfoHTTP(r, http.StatusNotFound)
			return
		}

		var notes []string
		for _, note := range folio.Notes {
			notes = append(notes, note.ListString())
		}
		jsonResponse, err := json.Marshal(notes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.ApplicationError(r, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
		logger.InfoHTTP(r, http.StatusOK)
	}).Methods("GET")

	// POST folios/{name} Append a note to a folio
	router.HandleFunc("/folios/{name}{slash:/?}", func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.ApplicationError(r, err)
			return
		}

		folio := folios[name]
		if folio == nil {
			w.WriteHeader(http.StatusNotFound)
			logger.InfoHTTP(r, http.StatusNotFound)
			return
		}

		err := folio.Append(r.FormValue("note"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.ApplicationError(r, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		logger.InfoHTTP(r, http.StatusCreated)
		logger.Info(fmt.Sprintf("Created note in folio %v\n", name))
	}).Methods("POST")

	// GET folios/{name}/{index}/done Toggle done on note
	router.HandleFunc("/folios/{name}/{index}/done{slash:/?}", func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		indexString := mux.Vars(r)["index"]

		index, err := strconv.Atoi(indexString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.ApplicationError(r, err)
			return
		}

		folio := folios[name]
		if folio == nil {
			w.WriteHeader(http.StatusNotFound)
			logger.InfoHTTP(r, http.StatusNotFound)
			return
		}

		if err = folio.ToggleDone(index); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Error(err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		logger.InfoHTTP(r, http.StatusCreated)
		logger.Info(fmt.Sprintf("Toggled done on note %v in folio %v\n", index, name))
	}).Methods("GET")

	// POST folios/ Create a folio
	router.HandleFunc("/folios{slash:/?}", func(w http.ResponseWriter, r *http.Request) {
		// Get folio name from request
		if err := r.ParseForm(); err != nil {
			logger.ApplicationError(r, err)
			return
		}
		name := r.FormValue("name")
		logger.Debug(name)

		// Validation
		matched, err := regexp.MatchString(`^[a-zA-Z]+$`, name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.ApplicationError(r, err)
			return
		}
		if !matched {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid folio name, must be one word")
			logger.InfoHTTP(r, http.StatusBadRequest)
			return
		}
		for _, f := range folios {
			if f.Name == name {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Folio with name exists, try a different name")
				logger.InfoHTTP(r, http.StatusBadRequest)
				return
			}
		}

		// Create New Folio
		folio, err := note.CreateFolio(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.ApplicationError(r, err)
			return
		}
		folios[name] = folio

		w.WriteHeader(http.StatusCreated)
		logger.InfoHTTP(r, http.StatusCreated)
		logger.Info(fmt.Sprintf("Created folio named %v\n", name))
	}).Methods("POST")

	// GET folios/ List all folio names
	router.HandleFunc("/folios{slash:/?}", func(w http.ResponseWriter, r *http.Request) {
		folioNames := []string{}
		for _, folio := range folios {
			folioNames = append(folioNames, folio.Name)
		}

		jsonResponse, err := json.Marshal(folioNames)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.ApplicationError(r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
		logger.InfoHTTP(r, http.StatusOK)
	}).Methods("GET")

	// DELETE folios/{name}
	router.HandleFunc("/folios/{name}{slash:/?}", func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]

		folio, ok := folios[name]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			logger.InfoHTTP(r, http.StatusNotFound)
			return
		}

		if err := folio.Delete(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.ApplicationError(r, err)
			return

		}

		delete(folios, name)

		w.WriteHeader(http.StatusOK)
		logger.InfoHTTP(r, http.StatusOK)
		logger.Info(fmt.Sprintf("Deleted folio %v\n", name))
	}).Methods("DELETE")

	// Manually reset 404 middleware or it will not fire. Custom 404 also ensures logging
	router.NotFoundHandler = router.NewRoute().HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		logger.InfoHTTP(r, http.StatusNotFound)
	}).GetHandler()

}

// Initializes Application Middleware
func initailizeMiddleware(router *mux.Router, logger *HTTPLogger.Logger) {
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
				w.WriteHeader(http.StatusUnauthorized)
				logger.InfoHTTP(r, http.StatusUnauthorized)
			}
		})
	})
}
