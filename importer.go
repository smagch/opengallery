package main

import (
	"encoding/json"
)

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

// func ImportExhibition(reader io.Reader) (err error) {
// 	// file, err := os.Open(filename)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// defer file.Close()
// 	// r := csv.NewReader(file)
// 	//var record []string
// 	var props
// 	r := csv.NewReader(reader)
// 	if props, err = r.Read(); err != nil {
// 		return
// 	}

// 	for {
// 		var record []string
// 		if record, err = r.Read(); err != nil && err != io.EOF {
// 			return
// 		}
// 		if record == nil {
// 			break
// 		}
// 		m := make(map[string]interface{})
// 		for i, prop := range props {
// 			if prop == "date_range" {
// 				date_range := strings.Split(record[i], ",")
// 				m[prop] = []string{
// 					strings.TrimSpace(date_range[0]),
// 					strings.TrimSpace(date_range[1]),
// 				}
// 			} else {
// 				m[prop] = record[i]
// 			}
// 		}

// 		var b []byte
// 		if b, err = json.Marshal(m); err != nil {
// 			return
// 		}
// 		e := &Exhibition{}
// 		if err = json.Unmarshal(b, &e); err != nil {
// 			return
// 		}
// 		if err = e.Create(); err != nil {
// 			return
// 		}
// 	}

// 	fmt.Println("ok")
// }
