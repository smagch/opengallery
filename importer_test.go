package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
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

func TestImportExhibition(t *testing.T) {
	b := []byte(`id,タイトル:title,説明:description,開始日:start,最終日:end
2014-1,新年おめでとう展【後期】,,2014/01/05,2014/01/13
2014-2,新春彫刻展,,2014/01/14,2014/01/20
2014-3,光彩画廊コレクション展,,2014/01/21,2014/01/27
2014-4,森清行・河原潤 二人展 二重星,,2014/01/28,2014/02/03`)

	galleryId := "B9FE1506-30C4-4CFF-B73E-99D859199A6D"
	reader := bytes.NewReader(b)
	exhibitions, err := ImportExhibition(galleryId, reader)
	if err != nil {
		t.Fatal(err)
	}
	desc := ""
	expected := []Exhibition{
		{"2014-1", galleryId, "新年おめでとう展【後期】", desc,
			*MustParseDateRange("2014-01-05", "2014-01-13")},
		{"2014-2", galleryId, "新春彫刻展", desc,
			*MustParseDateRange("2014-01-14", "2014-01-20")},
		{"2014-3", galleryId, "光彩画廊コレクション展", desc,
			*MustParseDateRange("2014-01-21", "2014-01-27")},
		{"2014-4", galleryId, "森清行・河原潤 二人展 二重星", desc,
			*MustParseDateRange("2014-01-28", "2014-02-03")},
	}

	for i, e := range expected {
		if !reflect.DeepEqual(e, exhibitions[i]) {
			t.Fatalf("Expected %v\n. But got %v instead", e, exhibitions[i])
		}
	}
}

func TestImportExhibitionNoContent(t *testing.T) {
	galleryId := "B9FE1506-30C4-4CFF-B73E-99D859199A6D"
	reader := bytes.NewReader([]byte{})
	if _, err := ImportExhibition(galleryId, reader); err != NoContentError {
		t.Fatal("It should return NoContentError")
	}
	reader = bytes.NewReader(
		[]byte(`id,タイトル:title,説明:description,開始日:start,最終日:end`),
	)
	if _, err := ImportExhibition(galleryId, reader); err != NoContentError {
		t.Fatal("It should return NoContentError")
	}
}

func TestImportFixture(t *testing.T) {
	if err := OpenTestDb(); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	MustTruncateAll()
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	filename := path.Join(pwd, "fixtures/hirama/hirama.json")
	err = ImportFixture(filename)
	if err != nil {
		t.Fatal(err)
	}
}
