package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func OpenTestDb() (err error) {
	db, err = sql.Open("postgres", "postgres://:@127.0.0.1/galleryinfo?sslmode=disable")
	db.SetMaxOpenConns(10)
	return
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func TestExhibitionMarshaling(t *testing.T) {
	b := []byte(`{
		"id": "2014-05-Foo",
		"gallery_id": "54b818d3-22f0-4f8b-6a04-170405fdb840",
		"title": "Foo",
		"description": "baar",
		"date_range": ["2014-05-10","2014-05-20"]
	}`)
	var m Exhibition
	var err error
	if err = json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	var b2 []byte
	if b2, err = json.Marshal(m); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	json.Compact(&buf, b)
	if buf.String() != string(b2[:]) {
		t.Fatalf("inconsistent json encoding/decoding.\n%v\n%v", buf.String(), string(b2[:]))
	}
}

func createRandomExhibition() *Exhibition {
	m := &Exhibition{}
	m.GalleryId = uuid.NewV4().String()
	m.Title = fmt.Sprintf("Exhibition-Title-%d", random(1000, 2000))
	m.Id = "ID:" + m.Title
	m.Description = "Description for " + m.Title
	dStart := time.Date(2014, time.Month(random(1, 12)), random(1, 29), 0, 0, 0, 0, time.UTC)
	dEnd := dStart.AddDate(0, 0, 14)
	m.DateRange = dateRange{dStart, dEnd}
	return m
}

func AssertSameExhibition(galleryId, id string, e *Exhibition) error {
	ex, err := GetExhibition(galleryId, id)
	if (err) != nil {
		return err
	}
	ex.DateRange[0] = ex.DateRange[0].In(time.UTC)
	ex.DateRange[1] = ex.DateRange[1].In(time.UTC)
	if !reflect.DeepEqual(e, ex) {
		return fmt.Errorf("Not deep equal\n%v\n\n%v", e, ex)
	}
	return nil
}

func SaveAndAssert(e *Exhibition, fn func() error) error {
	if err := fn(); err != nil {
		return err
	}
	if err := AssertSameExhibition(e.GalleryId, e.Id, e); err != nil {
		return err
	}
	return nil
}

func TestExhibitionCreate(t *testing.T) {
	if err := OpenTestDb(); err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	e := createRandomExhibition()
	if err := SaveAndAssert(e, e.Create); err != nil {
		t.Fatal(err)
	}

	e.Title = "Updated Title"
	e.DateRange[0] = e.DateRange[0].AddDate(0, -1, 0)
	e.DateRange[1] = e.DateRange[1].AddDate(0, 1, 0)
	e.Description = "Updated Description"
	if err := SaveAndAssert(e, e.Update); err != nil {
		t.Fatal(err)
	}
}
