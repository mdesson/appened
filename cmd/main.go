package main

import (
	"fmt"
	"net/http"

	"github.com/appened/note"
	"github.com/gorilla/mux"
)

// TODO: Add error handling on routes
// TODO: Add logging
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
	r.HandleFunc("/folios/{name}", getFolioHandler).Methods("GET")
	r.HandleFunc("/folios/{name}", appendNoteHandler).Methods("POST")
	r.HandleFunc("/folios", createFolioHandler).Methods("POST")

	http.ListenAndServe(":8081", r)
}

func createFolioHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Printf("Created folio named %v\n", r.FormValue("name"))
	w.WriteHeader(http.StatusCreated)
}

func getFolioHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	fmt.Println(w, "GET folios/%s", vars["name"])
	w.WriteHeader(http.StatusOK)
}

func appendNoteHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Printf("Created note with text %v\n", r.FormValue("note"))
	w.WriteHeader(http.StatusCreated)
}
