package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"math/rand"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"
)

func OpenTestDb() (err error) {
	db, err = sql.Open("postgres", "postgres://:@127.0.0.1/galleryinfo?sslmode=disable")
	db.SetMaxOpenConns(4)
	return
}

func MustParseDateRange(start, end string) *dateRange {
	d, err := parseDateRange(start, end)
	if err != nil {
		panic(err)
	}
	return d
}

func MustTruncateAll() {
	if _, err := db.Exec(`TRUNCATE exhibition, gallery`); err != nil {
		panic(err)
	}
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

type ExList []*VExhibition

func (e ExList) Len() int {
	return len(e)
}

func (e ExList) Less(i, j int) bool {
	return e[i].DateRange[1].Before(e[j].DateRange[1])
}

func (e ExList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
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

func GenerateRandomExhibition() *Exhibition {
	g := MustHaveGallery()
	m := &Exhibition{}
	m.GalleryId = g.Id
	m.Title = fmt.Sprintf("Exhibition-Title-%d", random(1000, 2000))
	m.Id = "ID:" + m.Title
	m.Description = "Description for " + m.Title
	dStart := time.Date(2014, time.Month(random(1, 12)), random(1, 29), 0, 0, 0, 0, time.UTC)
	dEnd := dStart.AddDate(0, 0, 14)
	m.DateRange = dateRange{dStart, dEnd}
	return m
}

func MustHaveExhibition() *Exhibition {
	e := GenerateRandomExhibition()
	if err := e.Create(); err != nil {
		panic(err)
	}
	return e
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

func insertExhibitionsWith(dr dateRange, gList []*Gallery) (results []Exhibition, err error) {
	var wg sync.WaitGroup
	wg.Add(len(gList))
	for _, g := range gList {
		e := Exhibition{}
		e.GalleryId = g.Id
		e.Title = fmt.Sprintf("Exhibition-Title-%d", random(1, 200000000))
		e.Id = "ID:" + e.Title
		e.Description = "Description for " + e.Title
		e.DateRange = dr
		results = append(results, e)
		go func(e *Exhibition) {
			err = e.Create()
			wg.Done()
		}(&e)
	}
	wg.Wait()
	return
}

func insertExhibitions(start time.Time, span int, total int) (eList []*Exhibition, err error) {
	// create gallery at first
	g := createRandomGallery()
	if err = g.Create(); err != nil {
		return
	}

	// create a long text
	addDescription := func(repeat int, e *Exhibition, txt string) {
		var b bytes.Buffer
		for i := 0; i < repeat; i++ {
			b.WriteString(txt)
		}
		e.Description = b.String()
	}

	var wg sync.WaitGroup
	var dateStart, dateEnd time.Time
	current := start

	for i := 0; i < total; i++ {
		wg.Add(1)
		e := &Exhibition{
			Id:        fmt.Sprintf("ID:%d", i),
			GalleryId: g.Id,
			Title:     fmt.Sprintf("Pagination-Test:%d", i),
		}
		addDescription(100, e, "This is the Description for "+e.Title)
		dateStart = current
		dateEnd = current.AddDate(0, 0, span)
		current = dateEnd.AddDate(0, 0, 1)
		e.DateRange = dateRange{dateStart, dateEnd}
		eList = append(eList, e)
		go func(i int, e *Exhibition, err error) {
			err = e.Create()
			wg.Done()
		}(i, e, err)
	}
	wg.Wait()
	return
}

func TestExhibitionCreate(t *testing.T) {
	if err := OpenTestDb(); err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	e := GenerateRandomExhibition()
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

	e.Description = "Updated Description"
	if err := SaveAndAssert(e, e.Sync); err != nil {
		t.Fatal(err)
	}

	e = GenerateRandomExhibition()
	if err := SaveAndAssert(e, e.Sync); err != nil {
		t.Fatal(err)
	}
}

func TestListExhibitionByGallery(t *testing.T) {
	if err := OpenTestDb(); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	MustTruncateAll()
	date := time.Date(2009, time.November, 10, 0, 0, 0, 0, time.UTC)
	eList, err := insertExhibitions(date, 7, 30)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = insertExhibitions(date, 7, 10); err != nil {
		t.Fatal(err)
	}

	var exhibitions []*VExhibition
	if exhibitions, err = ListExhibitionByGallery(eList[0].GalleryId); err != nil {
		t.Fatal(err)
	}

	if len(exhibitions) != 30 {
		t.Fatal("It should return 30 exhibitions")
	}

	for i, e := range exhibitions {
		ex := *eList[i]
		if ex.Title != e.Title || e.Gallery.Id != ex.GalleryId ||
			!ex.DateRange[0].Equal(e.DateRange[0]) ||
			!ex.DateRange[1].Equal(e.DateRange[1]) {
			t.Fatalf("Expected\n%#v\n. But got\n%#v", ex, e)
		}
	}
}

func TestGetExhibitionsByDateRange(t *testing.T) {
	if err := OpenTestDb(); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	MustTruncateAll()
	galleries, err := insertRandomGallery(20)
	if err != nil {
		t.Fatal(err)
	}
	span1 := MustParseDateRange("2014-01-15", "2014-01-20")
	span2 := MustParseDateRange("2014-01-19", "2014-01-21")

	if _, err = insertExhibitionsWith(*span1, galleries); err != nil {
		t.Fatal(err)
	}
	if _, err = insertExhibitionsWith(*span2, galleries); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		num   int
		start string
		end   string
	}{
		{0, "2013-01-14", "2013-01-14"},
		{0, "2014-01-14", "2014-01-14"},
		{20, "2014-01-15", "2014-01-15"},
		{20, "2014-01-16", "2014-01-16"},
		{20, "2014-01-17", "2014-01-17"},
		{20, "2014-01-18", "2014-01-18"},
		{40, "2014-01-19", "2014-01-19"},
		{40, "2014-01-20", "2014-01-20"},
		{20, "2014-01-21", "2014-01-21"},
		{20, "2014-01-21", "2014-01-21"},
		{0, "2014-01-22", "2014-01-22"},
		{0, "2015-01-21", "2015-01-21"},
		{40, "2014-01-14", "2014-01-22"},
	}

	for _, c := range cases {
		dr := MustParseDateRange(c.start, c.end)
		results, err := SearchExhibitions(dr)
		if err != nil {
			t.Fatal(err)
		}
		num := len(results)
		if num != c.num {
			t.Fatalf("Expected %d results length on %s. But got %d", c.num,
				c.start, num)
		}
		if !sort.IsSorted(ExList(results)) {
			t.Fatal("Response should be sorted in the final date order")
		}
	}
}
