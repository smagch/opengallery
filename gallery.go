package main

import (
	"encoding/json"
	"fmt"
)

// Gallery represents gallery model.
type Gallery struct {
	Id    string          `json:"id"`
	Name  string          `json:"name"`
	Meta  json.RawMessage `json:"meta"`
	About string          `json:"about"`
}

// Validate returns error if a field value is invalid.
func (g *Gallery) Validate() error {
	if !IsUUID(g.Id) {
		return fmt.Errorf("Invalid Id: %s is not an UUID", g.Id)
	}
	return nil
}

// Create insert a row in gallery table.
func (g *Gallery) Create() error {
	if err := g.Validate(); err != nil {
		return err
	}
	_, err := db.Exec(`
		INSERT INTO
			gallery (id, name, meta, about)
		VALUES
			($1, $2, $3, $4)`,
		g.Id, g.Name, string(g.Meta), g.About)
	return err
}

// GetGallery fetch a row from gallry table.
func GetGallery(id string) (g *Gallery, err error) {
	g = &Gallery{}
	err = db.QueryRow(`
		SELECT
			id, name, meta, about
		FROM
			gallery
		WHERE
			id = $1`,
		id).Scan(&g.Id, &g.Name, &g.Meta, &g.About)
	return
}