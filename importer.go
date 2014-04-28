package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
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
