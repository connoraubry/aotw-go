package server

import (
	"fmt"
	"time"
)

type Index struct {
	Album AlbumInfo
	Form  Form
}

type AlbumInfo struct {
	Title  string
	Artist string
	Year   int

	SubmittedOn time.Time
	SubmittedBy string

	Week     int
	ChosenOn time.Time
	Image    string
}

type Form struct {
	Result  string
	Options Options
}

type Options []FormOption
type FormOption struct {
	Title  string
	Artist string
}

func (self *AlbumInfo) GenInsertQuery(database string) string {

	cols := `
	(title, artist, date, submitted_on,
	submitted_by, week, chosen_on, image)`

	vals := fmt.Sprintf(
		"(%s, %s, %d, %d, %s, %d, %d, %s)",
		self.Title,
		self.Artist,
		self.Year,
		self.SubmittedOn.Unix(),
		self.SubmittedBy,
		self.Week,
		self.ChosenOn.Unix(),
		self.Image,
	)

	query := fmt.Sprintf("INSERT INTO %s %s VALUES%s", database, cols, vals)

	return query
}
