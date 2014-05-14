package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"reflect"
	"strconv"
	"sync"
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

func MustHaveGallery() *Gallery {
	g := createRandomGallery()
	if err := g.Create(); err != nil {
		panic(err)
	}
	return g
}

func insertRandomGallery(total int) (results []*Gallery, err error) {
	var wg sync.WaitGroup
	wg.Add(total)

	for i := 0; i < total; i++ {
		g := createRandomGallery()
		results = append(results, g)
		go func(g *Gallery, err error) {
			err = g.Create()
			wg.Done()
		}(g, err)
	}

	wg.Wait()
	return
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
	if err := g.Create(); err != nil {
		t.Fatal(err)
	}
	AssertSameGallery(g.Id, g)

	g.Name = "Updated Gallery Name:" + g.Id
	g.About = "Updated About:" + g.Id
	g.Meta = []byte(`{"location":"updated"}`)
	if err := g.Update(); err != nil {
		t.Fatal(err)
	}
	AssertSameGallery(g.Id, g)
}

func TestSyncGallery(t *testing.T) {
	if err := OpenTestDb(); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	g := createRandomGallery()
	if err := g.Sync(); err != nil {
		t.Fatal(err)
	}
	AssertSameGallery(g.Id, g)

	g.Name = "Updated Gallery Name:" + g.Id
	g.About = "Updated About:" + g.Id
	g.Meta = []byte(`{"location":"updated"}`)
	if err := g.Sync(); err != nil {
		t.Fatal(err)
	}
	AssertSameGallery(g.Id, g)
}
