package main

import (
	"encoding/json"
	"github.com/smagch/patree"
	"io"
	"net/http"
	"strconv"
)

var (
	Status500 = []byte(`{"message": "InternalServerError"}`)
)

func New404(urlStr string) *NotFoundError {
	return &NotFoundError{urlStr}
}

type NotFoundError struct {
	url string
}

func (err *NotFoundError) Error() string {
	return "URL " + err.url + "NotFound"
}

type ListResponse struct {
	Results interface{} `json:"results,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// HandleError
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if v, ok := err.(*ValidationError); ok {
		BadRequest(w, v)
	} else if _, ok := err.(*NotFoundError); ok {
		NotFound(w)
	} else {
		InternalServerError(w, err)
	}
}

// BadRequest
// TODO respond JSON instead
func BadRequest(w http.ResponseWriter, err *ValidationError) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, err.Error())
}

// InternalServerError
// TODO logging
func InternalServerError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, err.Error())
}

// NotFound sends 404 notfound status.
func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func Boot(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return nil
}

func Json(w http.ResponseWriter, model interface{}) {
	b, err := json.Marshal(model)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}

// App returns main http multiplexer.
func App() *patree.PatternTreeServeMux {
	mux := patree.New()
	mux.UseFunc(Boot)
	mux.Error(HandleError)

	exHandler := &ExhibitionHandler{"exhibition_id", "gallery_id", "date"}
	mux.Get("/galleries/<uuid:gallery_id>/exhibitions/<exhibition_id>",
		exHandler.Get)
	mux.Get("/galleries/<uuid:gallery_id>/exhibitions", exHandler.ListByGallery)
	mux.Get("/exhibitions/<date:date>", exHandler.FindByDate)

	gHandler := &GalleryHandler{"gallery_id"}
	mux.Get("/galleries/<uuid:gallery_id>", gHandler.Get)
	return mux
}
