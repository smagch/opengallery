package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type routeTest struct {
	urlFormat string
	cases     []routeCase
}

type routeCase struct {
	urlArgs []string
	code    int
	v       responseValidator
}

type responseValidator interface {
	validate(b []byte) error
}

// ListResponse
type listResponse struct {
	Results json.RawMessage `json:"results,omitempty"`
	Errors  json.RawMessage `json:"errors,omitempty"`
	Code    int             `json:"code,omitempty"`
}

type byteResponseValidator struct {
	b []byte
}

func (v *byteResponseValidator) validate(b []byte) error {
	if !bytes.Equal(b, v.b) {
		return fmt.Errorf("Invalid Response: %s\nExpected %s\n", string(b),
			string(v.b))
	}
	return nil
}

type exhibitionResponseValidator struct {
	data []Exhibition
}

func (v *exhibitionResponseValidator) validate(b []byte) error {
	var res listResponse
	err := json.Unmarshal(b, &res)
	if err != nil {
		return err
	}
	var exhibitions []Exhibition
	err = json.Unmarshal(res.Results, &exhibitions)
	if err != nil {
		return err
	}
	if len(v.data) != len(exhibitions) {
		return fmt.Errorf("Inconsistent response. Expected %v. Got %v instead",
			v.data, exhibitions)
	}
	if !reflect.DeepEqual(v.data, exhibitions) {
		return errors.New("Inconsistent exhibition list")
	}
	return nil
}

func (c routeCase) getArgs() []interface{} {
	args := make([]interface{}, len(c.urlArgs))
	for i, arg := range c.urlArgs {
		args[i] = interface{}(arg)
	}
	return args
}

func (rt *routeTest) exec(t *testing.T) {
	mux := App()
	for _, c := range rt.cases {
		args := c.getArgs()
		urlStr := fmt.Sprintf(rt.urlFormat, args...)
		r, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if c.code != w.Code {
			t.Errorf("Status code should be %d rather than %d with %s", c.code,
				w.Code, urlStr)
			if w.Code == 500 {
				t.Fatal(w.Body.String())
			}
		}
		if c.code < 300 && c.v != nil {
			c.v.validate(w.Body.Bytes())
		}
	}
}

func TestExhibitionRoutes(t *testing.T) {
	if err := OpenTestDb(); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	MustTruncateAll()
	e := MustHaveExhibition()
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	v := &byteResponseValidator{b}
	rt := &routeTest{"/galleries/%s/exhibitions/%s", []routeCase{
		{[]string{strings.ToUpper(e.GalleryId), e.Id}, 200, v},
		{[]string{strings.ToLower(e.GalleryId), e.Id}, 200, v},
		{[]string{e.GalleryId, "invalid-id"}, 404, nil},
		{[]string{"invalid-gallery-id", "invalid-id"}, 404, nil},
		{[]string{"invalid-gallery-id", e.Id}, 404, nil},
		{[]string{uuid.NewV4().String(), e.Id}, 404, nil},
	}}
	rt.exec(t)
}

func TestExhibitionByDate(t *testing.T) {
	if err := OpenTestDb(); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	MustTruncateAll()

	galleries, err := insertRandomGallery(20)
	if err != nil {
		t.Fatal(err)
	}

	span := MustParseDateRange("2014-04-02", "2014-04-04")
	var exhibitions []Exhibition
	exhibitions, err = insertExhibitionsWith(*span, galleries)
	if err != nil {
		t.Fatal(err)
	}

	emptyResults := ListResponse{Results: []interface{}{}}
	empty, err := json.Marshal(emptyResults)
	if err != nil {
		t.Fatal(err)
	}
	emptyV := &byteResponseValidator{empty}
	v := &exhibitionResponseValidator{exhibitions}

	rt := &routeTest{"/exhibitions/%s", []routeCase{
		{[]string{"2014-04-01"}, 200, emptyV},
		{[]string{"2014-04-02"}, 200, v},
		{[]string{"2014-04-03"}, 200, v},
		{[]string{"2014-04-04"}, 200, v},
		{[]string{"2014-04-05"}, 200, emptyV},
	}}

	rt.exec(t)
}

func TestGalleryRoutes(t *testing.T) {
	if err := OpenTestDb(); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	MustTruncateAll()
	g := MustHaveGallery()
	b, err := json.Marshal(g)
	if err != nil {
		t.Fatal(err)
	}
	v := &byteResponseValidator{b}
	rt := &routeTest{"/galleries/%s", []routeCase{
		{[]string{strings.ToUpper(g.Id)}, 200, v},
		{[]string{strings.ToLower(g.Id)}, 200, v},
		{[]string{"hogehoge"}, 404, nil},
		{[]string{g.Id + "a"}, 404, nil},
		{[]string{uuid.NewV4().String()}, 404, nil},
	}}
	rt.exec(t)
}
