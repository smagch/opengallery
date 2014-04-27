package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseGalleryData(t *testing.T) {
	b := []byte(`{
		"id": "B9FE1506-30C4-4CFF-B73E-99D859199A6D",
		"name": "ヒラマ画廊",
		"address": "070-0032 旭川市２条通８丁目",
		"about": "Test",
		"open_at": "10:00",
		"close_at": "18:00",
		"close_on": "",
		"exhibitions": [
			"2013.csv",
			"2014.csv"
		]
	}`)
	g, exhibitions, err := ParseGalleryData(b)
	if err != nil {
		t.Fatal(err)
	}

	gExpected := &Gallery{
		Id:    "B9FE1506-30C4-4CFF-B73E-99D859199A6D",
		Name:  "ヒラマ画廊",
		About: "Test",
		Meta:  []byte{},
	}
	exExpected := []string{"2013.csv", "2014.csv"}
	metaExpected := map[string]string{
		"open_at":  "10:00",
		"close_at": "18:00",
		"close_on": "",
		"address":  "070-0032 旭川市２条通８丁目",
	}
	meta := map[string]string{}
	if err := json.Unmarshal(g.Meta, &meta); err != nil {
		t.Fatal(err)
	}
	g.Meta = []byte{}

	if !reflect.DeepEqual(gExpected, g) {
		t.Fatalf("Expected: %v\nGot %v instead\n", gExpected, g)
	}
	if !reflect.DeepEqual(exExpected, exhibitions) {
		t.Fatalf("Expected %v\nGot %v instead\n", exExpected, exhibitions)
	}
	if !reflect.DeepEqual(metaExpected, meta) {
		t.Fatalf("Expected: %v\nGot %v instead\n", metaExpected, meta)
	}
}
