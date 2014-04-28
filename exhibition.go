package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"strings"
	"time"
)

const (
	DATE_LAYOUT       = "2006-01-02"
	DATE_LAYOUT_SLASH = "2006/01/02"
)

var (
	db *sql.DB
)

type dateRange [2]time.Time

func (dr dateRange) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`["%s","%s"]`, dr[0].Format(DATE_LAYOUT),
		dr[1].Format(DATE_LAYOUT))), nil
}

// len(`["2014-05-10","2014-05-20"]`) == 27
func (dr *dateRange) UnmarshalJSON(data []byte) error {
	l := len(data)
	if l < 27 {
		return errors.New("DateRange Parse Error: not enough length")
	}
	if data[0] != '[' || data[l-1] != ']' {
		return errors.New("DateRange Parse Error: daterange should be an array")
	}
	d := strings.Split(string(data[1:l-1]), ",")
	if len(d) != 2 {
		return errors.New("DateRange Parse Error: DataRange should have two item")
	}
	dr2, err := parseDateRange(strings.Trim(d[0], `" `), strings.Trim(d[1], `" `))
	if err != nil {
		return err
	}
	dr[0], dr[1] = dr2[0], dr2[1]
	return nil
}

func (dr *dateRange) Format() string {
	return fmt.Sprintf("[%s,%s]", dr[0].Format(DATE_LAYOUT),
		dr[1].Format(DATE_LAYOUT))
}

func parseDateRange(start, end string) (*dateRange, error) {
	return parseDateRangeByLayout(start, end, DATE_LAYOUT)
}

func ParseDateRangeBySlash(start, end string) (dr *dateRange, err error) {
	return parseDateRangeByLayout(start, end, DATE_LAYOUT_SLASH)
}

func parseDateRangeByLayout(start, end, layout string) (dr *dateRange, err error) {
	var dStart, dEnd time.Time
	if dStart, err = time.Parse(layout, start); err != nil {
		return nil, errors.New("DateRange Parse Error: Invalid Date Start " + start)
	}
	if dEnd, err = time.Parse(layout, end); err != nil {
		return nil, errors.New("DateRange Parse Error: Invalid Date End " + end)
	}
	return &dateRange{dStart, dEnd}, nil
}

//
type Exhibition struct {
	Id          string    `json:"id"`
	GalleryId   string    `json:"gallery_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DateRange   dateRange `json:"date_range"`
}

func (e *Exhibition) GetByteId() []byte {
	b := [32]byte{}
	date := e.GetDateByte()
	hash := e.GetHashId()
	copy(b[0:], date)
	copy(b[4:], hash)
	return b[:]
}

func (e *Exhibition) GetHashId() []byte {
	h := sha256.New224()
	io.WriteString(h, strings.ToLower(e.GalleryId))
	io.WriteString(h, e.Id)
	return h.Sum(nil)
}

func (e *Exhibition) GetDateByte() []byte {
	t := e.DateRange[0]
	n := t.Day() + int(t.Month())*32 + t.Year()*32*16
	u := make([]byte, 4)
	binary.BigEndian.PutUint32(u, uint32(n))
	return u
}

// Validate executes validation for exhibition properties.
func (e *Exhibition) Validate() (err ValidationError) {
	if len(e.Id) == 0 {
		err = err.Append("Invalid id: " + e.Id + " should not empty")
	}
	if !IsUUID(e.GalleryId) {
		err = err.Append("Invalid gallery_id: " + e.GalleryId + " should be UUID")
	}
	return
}

// Create insert a row into exhibition table.
func (e *Exhibition) Create() error {
	if err := e.Validate(); err != nil {
		return err
	}
	b := e.GetByteId()
	_, err := db.Exec(`
		INSERT INTO
			exhibition
			(id, _byteid, gallery_id, title, description, date_range)
		VALUES
			($1, $2, $3, $4, $5, $6)
	`, e.Id, b, e.GalleryId, e.Title, e.Description, e.DateRange.Format())
	return err
}

// Update update an exhibition row
func (e *Exhibition) Update() error {
	if err := e.Validate(); err != nil {
		return err
	}
	b := e.GetByteId()
	hashId := e.GetHashId()
	_, err := db.Exec(`
		UPDATE
			exhibition
		SET
			(_byteid, title, description, date_range) = ($2, $3, $4, $5)
		WHERE
			substring(_byteid, 5) = $1
		`, hashId, b, e.Title, e.Description, e.DateRange.Format())
	return err
}

// CreateOrUpdate update if exists. If not create new model.
func (e *Exhibition) CreateOrUpdate() error {
	if err := e.Validate(); err != nil {
		return err
	}
	var exists bool
	hashId := e.GetHashId()
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM exhibition WHERE substring(_byteid, 5) = $1
		)
	`, hashId).Scan(&exists)

	if err != nil {
		return err
	}
	if exists {
		err = e.Update()
	} else {
		err = e.Create()
	}
	return err
}

// GetExhibition fetch an exhibition model.
func GetExhibition(galleryId, id string) (*Exhibition, error) {
	var dateStart, dateEnd time.Time
	e := &Exhibition{
		GalleryId: galleryId,
		Id:        id,
	}
	if err := e.Validate(); err != nil {
		return nil, err
	}
	// use lower case uuid for consistency
	e.GalleryId = strings.ToLower(e.GalleryId)

	b := e.GetHashId()
	err := db.QueryRow(`
		SELECT
			title, description, lower(date_range), upper(date_range)
		FROM
			exhibition
		WHERE
			substring(_byteid, 5) = $1
		`, b).Scan(&e.Title, &e.Description, &dateStart, &dateEnd)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	dateEnd = dateEnd.AddDate(0, 0, -1)
	e.DateRange = dateRange{dateStart, dateEnd}
	return e, nil
}

func GetExhibitionJSON(galleryId, id string) (b []byte, err error) {
	var e *Exhibition
	if e, err = GetExhibition(galleryId, id); err != nil || e == nil {
		return
	}
	b, err = json.Marshal(*e)
	return
}

func SearchExhibitions(dr *dateRange) (results []Exhibition, err error) {
	var rows *sql.Rows
	rows, err = db.Query(`
		SELECT
			gallery_id, title, lower(date_range), upper(date_range)
		FROM
			exhibition
		WHERE
			date_range && $1
		ORDER BY
			upper(date_range)
		LIMIT 100
	`, dr.Format())

	if err != nil {
		return
	}
	defer rows.Close()

	results = []Exhibition{}
	for rows.Next() {
		var start, end time.Time
		e := Exhibition{}
		err = rows.Scan(&e.GalleryId, &e.Title, &start, &end)
		if err != nil {
			return nil, err
		}
		end = end.AddDate(0, 0, -1)
		e.DateRange = dateRange{start, end}
		results = append(results, e)
	}

	err = rows.Err()
	return
}
