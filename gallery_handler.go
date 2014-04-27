package main

import (
	"github.com/smagch/patree"
	"net/http"
)

// Gallery
type GalleryHandler struct {
	IdName string
}

func (h *GalleryHandler) Get(w http.ResponseWriter, r *http.Request) error {
	id := patree.Param(r, h.IdName)
	g, err := GetGallery(id)
	if err != nil {
		return err
	} else if g == nil {
		return New404(r.URL.Path)
	}
	Json(w, g)
	return nil
}
