package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"net/http/httptest"
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
	body    []byte
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
		if c.code < 300 && !bytes.Equal(c.body, w.Body.Bytes()) {
			t.Fatalf("Invalid Response: %s\nExpected %s\n", w.Body.String(),
				string(c.body))
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
	rt := &routeTest{"/galleries/%s/exhibitions/%s", []routeCase{
		{[]string{strings.ToUpper(e.GalleryId), e.Id}, 200, b},
		{[]string{strings.ToLower(e.GalleryId), e.Id}, 200, b},
		{[]string{e.GalleryId, "invalid-id"}, 404, nil},
		{[]string{"invalid-gallery-id", "invalid-id"}, 404, nil},
		{[]string{"invalid-gallery-id", e.Id}, 404, nil},
		{[]string{uuid.NewV4().String(), e.Id}, 404, nil},
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
	rt := &routeTest{"/galleries/%s", []routeCase{
		{[]string{strings.ToUpper(g.Id)}, 200, b},
		{[]string{strings.ToLower(g.Id)}, 200, b},
		{[]string{"hogehoge"}, 404, nil},
		{[]string{g.Id + "a"}, 404, nil},
		{[]string{uuid.NewV4().String()}, 404, nil},
	}}
	rt.exec(t)
}
