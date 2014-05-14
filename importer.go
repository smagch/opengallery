package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var NoContentError = errors.New("No Content")

type galleryInput struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	About       string   `json:"about"`
	Address     string   `json:"address"`
	OpenAt      string   `json:"open_at"`
	CloseAt     string   `json:"close_at"`
	CloseOn     string   `json:"close_on"`
	Exhibitions []string `json:"exhibitions"`
}

// TODO log unknown attributes
func ParseGalleryData(b []byte) (g *Gallery, exhibitions []string, err error) {
	input := &galleryInput{}
	if err = json.Unmarshal(b, input); err != nil {
		return
	}

	meta := map[string]string{
		"address":  input.Address,
		"open_at":  input.OpenAt,
		"close_at": input.CloseAt,
		"close_on": input.CloseOn,
	}

	g = &Gallery{
		Id:    input.Id,
		Name:  input.Name,
		About: input.About,
	}
	if g.Meta, err = json.Marshal(meta); err != nil {
		return
	}

	exhibitions = input.Exhibitions
	return
}

func ImportExhibition(galleryId string, reader io.Reader) (exhibitions []Exhibition, err error) {
	var props []string
	r := csv.NewReader(reader)
	if props, err = r.Read(); err != nil {
		if err == io.EOF {
			err = NoContentError
			return
		}
		return
	}

	propsRequired := []string{"id", "title", "description", "start", "end"}
	propsOptional := []string{"alerts", "notes"}
	propsAllowed := append(propsRequired, propsOptional...)
	usedProps := make(map[int]bool)

	for _, p := range propsAllowed {
		for i, verboseProp := range props {
			if strings.HasSuffix(verboseProp, p) {
				props[i] = p
				usedProps[i] = true
				continue
			}
		}
		// if p != "alerts" && p != "notes" {
		// 	return nil, fmt.Errorf("property \"%s\" is required.", p)
		// }
	}

	exhibitions = []Exhibition{}
	for {
		var record []string
		if record, err = r.Read(); err != nil {
			if err == io.EOF {
				if len(exhibitions) == 0 {
					err = NoContentError
				} else {
					err = nil
				}
			}
			return
		}

		if record == nil {
			break
		}
		m := make(map[string]interface{})
		var dateStart, dateEnd string
		for i, prop := range props {
			if _, ok := usedProps[i]; !ok {
				continue
			}
			if prop == "start" {
				dateStart = record[i]
			} else if prop == "end" {
				dateEnd = record[i]
			} else if prop != "alerts" && prop != "notes" {
				// TODO alerts and notes
				m[prop] = record[i]
			}
		}
		m["date_range"], err = ParseDateRangeBySlash(dateStart, dateEnd)
		if err != nil {
			return
		}

		var b []byte
		if b, err = json.Marshal(m); err != nil {
			return
		}
		e := Exhibition{}
		if err = json.Unmarshal(b, &e); err != nil {
			return
		}
		e.GalleryId = galleryId
		exhibitions = append(exhibitions, e)
	}

	return
}

// 1.   make an http request to the url
// 2.0  Parse and validate the JSON data
// 2.1  Trim and create a checksum for the data.
// 2.2  make HEAD requests to the exhibitions url
// 2.3  return error if exhibition url does not exist.
// 3.0  Create or update gallery data.
// 3.1  Create or update exhibition data.
// func ImportGallery(url string) error {}

func exists(name string) (bool, error) {
	fi, err := os.Stat(name)
	if err == nil {
		return !fi.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// ImportFixture imports data from the given filename. This is a ad-hoc
// implementation that would be replaced http loading in the future.
func ImportFixture(filename string) error {
	if path.Ext(filename) != ".json" {
		return errors.New("Only JSON files are supported")
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	var b []byte
	if b, err = ioutil.ReadAll(file); err != nil {
		return err
	}

	var g *Gallery
	var exhibitions []string
	if g, exhibitions, err = ParseGalleryData(b); err != nil {
		return err
	}

	if vError := g.Validate(); vError != nil {
		return vError
	}

	// check existance
	dirname := path.Dir(filename)
	for i, name := range exhibitions {
		exhibitions[i] = path.Join(dirname, name)
		var ok bool
		ok, err = exists(exhibitions[i])
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("No such file as %s. File %s does not exists",
				exhibitions[i], name)
		}
	}

	if err = g.Sync(); err != nil {
		return err
	}

	for _, filename := range exhibitions {
		var f *os.File
		if f, err = os.Open(filename); err != nil {
			return err
		}
		var exList []Exhibition
		if exList, err = ImportExhibition(g.Id, f); err != nil {
			return err
		}
		for _, e := range exList {
			if err := e.Sync(); err != nil {
				return err
			}
		}
	}
	return nil
}
