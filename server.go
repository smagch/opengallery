package main

import (
	"github.com/smagch/patree"
	"io"
	"net/http"
	"strconv"
)

var (
	Status500 = []byte(`{"message": "InternalServerError"}`)
)

// HandleError
func HandleError(w http.ResponseWriter, err error) {
	if v, ok := err.(*ValidationError); ok {
		BadRequest(w, v)
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

// ExhibitionHandler handles exhibition resources.
type ExhibitionHandler struct {
	IdName        string
	GalleryIdName string
}

// Get send a JSON response that represents an exhibition.
func (h *ExhibitionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, galleryId := patree.Param(r, h.IdName), patree.Param(r, h.GalleryIdName)
	var err error
	var b []byte
	// TODO it's not quite suitable to respond bad request with validation error
	if b, err = GetExhibitionJSON(galleryId, id); err != nil {
		HandleError(w, err)
		return
	}
	if b == nil {
		NotFound(w)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}

// App returns main http multiplexer.
func App() *patree.PatternTreeServeMux {
	mux := patree.New()
	h := &ExhibitionHandler{"exhibition_id", "gallery_id"}
	mux.Get("/galleries/<uuid:gallery_id>/exhibitions/<exhibition_id>",
		http.HandlerFunc(h.Get))
	return mux
}
