package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/tui95/go-url-shortener/internal/database"
	"github.com/tui95/go-url-shortener/internal/lib"
)

type Router struct {
	*http.ServeMux
	db *sql.DB
}

func NewRouter(db *sql.DB) *Router {
	mux := http.NewServeMux()
	router := Router{ServeMux: mux, db: db}
	router.HandleFunc("POST /", router.CreateShortUrlHandler)
	router.HandleFunc("GET /{base64Id}", router.RedirectToOriginalUrlHandler)
	return &router
}

type UrlShorten struct {
	Url string `json:"url"`
}

func (router *Router) CreateShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody UrlShorten
	err := parseBody(w, r, &requestBody)
	if err != nil {
		return
	}

	id, err := database.CreateUrlMapping(router.db, requestBody.Url)
	if err != nil {
		internalServerError(w, err)
		return
	}

	shortUrl, err := url.JoinPath(getBaseUrl(r), lib.Encode(id))
	if err != nil {
		internalServerError(w, err)
		return
	}
	responseData := UrlShorten{Url: shortUrl}
	err = writeJsonResponse(w, http.StatusCreated, responseData)
	if err != nil {
		log.Printf("Unable to encode response: %v\n", err)
		internalServerError(w, err)
	}
}

type ErrorResponse struct {
	Detail string `json:"detail"`
}

func (router *Router) RedirectToOriginalUrlHandler(w http.ResponseWriter, r *http.Request) {
	base64Id := r.PathValue("base64Id")
	id := lib.Decode(base64Id)
	url, err := database.GetUrlById(router.db, id)
	if errors.Is(err, sql.ErrNoRows) {
		notFound(w, err)
		return
	} else if err != nil {
		internalServerError(w, err)
		return
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func parseBody[T any](w http.ResponseWriter, r *http.Request, t *T) error {
	err := json.NewDecoder(r.Body).Decode(&t)
	if errors.Is(err, io.EOF) {
		badRequest(w, nil, ErrorResponse{Detail: "Required body."})
		return err
	} else if err != nil {
		log.Printf("Unable to parse body: %+v\n", err)
		internalServerError(w, err)
		return err
	}
	return nil
}

func getBaseUrl(r *http.Request) string {
	var protocol string
	if r.TLS != nil {
		protocol = "https"
	} else {
		protocol = "http"
	}
	return fmt.Sprintf("%v://%v", protocol, r.Host)
}

func internalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, err error, response any) {
	err = writeJsonResponse(w, http.StatusBadRequest, response)
	if err != nil {
		internalServerError(w, err)
	}
}

func notFound(w http.ResponseWriter, err error) {
	err = writeJsonResponse(w, http.StatusNotFound, ErrorResponse{Detail: "Not found."})
	if err != nil {
		internalServerError(w, err)
	}
}

func writeJsonResponse(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}
