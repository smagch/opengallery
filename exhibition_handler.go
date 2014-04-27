package main

import (
	"github.com/smagch/patree"
	"net/http"
	"time"
)

// ExhibitionHandler handles exhibition resources.
type ExhibitionHandler struct {
	IdName        string
	GalleryIdName string
	DateName      string
}

// Get send a JSON response that represents an exhibition.
func (h *ExhibitionHandler) Get(w http.ResponseWriter, r *http.Request) error {
	id, galleryId := patree.Param(r, h.IdName), patree.Param(r, h.GalleryIdName)
	// TODO it's not quite suitable to respond bad request with validation error
	e, err := GetExhibition(galleryId, id)
	if err != nil {
		return err
	} else if e == nil {
		return New404(r.URL.Path)
	}
	Json(w, e)
	return nil
}

func (h *ExhibitionHandler) FindByDate(w http.ResponseWriter, r *http.Request) error {
	date := patree.Param(r, h.DateName)
	d, err := time.Parse(DATE_LAYOUT, date)
	if err != nil {
		return err
	}
	dr := &dateRange{d, d}
	var results []Exhibition
	results, err = SearchExhibitions(dr)
	if err != nil {
		return err
	}
	res := &ListResponse{Results: results}
	Json(w, res)
	return nil
}
