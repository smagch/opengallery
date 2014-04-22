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

	mux := App()

	cases := []struct {
		galleryId string
		id        string
		code      int
		body      []byte
	}{
		{strings.ToUpper(e.GalleryId), e.Id, 200, b},
		{strings.ToLower(e.GalleryId), e.Id, 200, b},
		{e.GalleryId, "invalid-id", 404, nil},
		{"invalid-gallery-id", "invalid-id", 404, nil},
		{"invalid-gallery-id", e.Id, 404, nil},
		{uuid.NewV4().String(), e.Id, 404, nil},
	}

	for _, c := range cases {
		urlStr := fmt.Sprintf("/galleries/%s/exhibitions/%s", c.galleryId, c.id)
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
		if c.code == 200 && !bytes.Equal(b, w.Body.Bytes()) {
			t.Fatalf("Invalid Response: %s\nExpected %s\n", w.Body.String(),
				string(b))
		}
	}
}
