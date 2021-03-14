package server

import (
	"encoding/json"
	"fmt"
	"github.com/seungha-kim/urlo-go/internal/pkg/conf"
	"github.com/seungha-kim/urlo-go/internal/pkg/util"
	"log"
	"net/http"
)

func passedAccessKeyFilter(w http.ResponseWriter, req *http.Request) bool {
	accessKeyParam := req.FormValue("access_key")
	if accessKeyParam != conf.AccessKey {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, "401 Unauthorized")
		return false
	}
	return true
}

func handleUnexpectedErr(w http.ResponseWriter, err error) {
	log.Println("Unexpected error -", err)
	w.WriteHeader(500)
	_, _ = fmt.Fprintf(w, "500 Internal Server Error")
}

func authHandler(w http.ResponseWriter, req *http.Request) {
	if !passedAccessKeyFilter(w, req) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func createHandler(w http.ResponseWriter, req *http.Request) {
	if !passedAccessKeyFilter(w, req) {
		return
	}

	url := req.FormValue("url")
	id, err := CreateIdByUrl(db, url)
	if err != nil {
		handleUnexpectedErr(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "id": id})
}

func redirectHandler(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("u")
	url, err := GetUrlById(db, id)
	if err != nil {
		handleUnexpectedErr(w, err)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(301)
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == "GET" && len(req.FormValue("u")) > 0:
		redirectHandler(w, req)
	case req.Method == "POST":
		createHandler(w, req)
	}
}

func RunServer() {
	log.Println("Starting server at", conf.ServerAddress)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/auth", authHandler)
	util.PanicIfError(http.ListenAndServe(conf.ServerAddress, nil), "Failed to start server")
}
