package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"reflect"
	"strconv"
	"testing"
)

func createRandomGallery() *Gallery {
	g := &Gallery{}
	g.Id = uuid.NewV4().String()
	i := strconv.Itoa(random(1, 10000))
	g.Name = "Gallery:" + i
	g.About = "AboutMe:" + i
	g.Meta = []byte(`{"location": "Asahikawa City"}`)
	return g
}

func AssertSameGallery(id string, g *Gallery) error {
	gallery, err := GetGallery(id)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(g, gallery) {
		return fmt.Errorf("Not deep equal\n%v\n\n%v", g, gallery)
	}
	return nil
}

func TestGalleryMarshaling(t *testing.T) {
	b := []byte(`{
		"id": "9fc312ff-2d94-47dd-a643-c69c63294624",
		"name": "Foobar",
		"meta": {"location":"Asahikawa City"},
		"about": "About me"
	}`)

	var g *Gallery
	var err error
	if err = json.Unmarshal(b, &g); err != nil {
		t.Fatal(err)
	}
	var b2 []byte
	if b2, err = json.Marshal(g); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	json.Compact(&buf, b2)
	if buf.String() != string(b2[:]) {
		t.Fatal("inconsistent json encoding/decoding")
	}
}

func TestCreateGallery(t *testing.T) {
	if err := OpenTestDb(); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	g := createRandomGallery()
	var id string
	if err := g.Create(); err != nil {
		t.Fatal(err)
	}

	AssertSameGallery(id, g)
}