package server

import (
	"database/sql"
	"os"

	"github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Path string
}

func NewDatabase(path string) *Database {
	DBStruct := Database{Path: path}

	_, err := os.Stat(path)
	if err != nil {
	}

	db := DBStruct.loadDB()
	defer db.Close()
	DBStruct.initDB(db)

	return &DBStruct

}

func (self *Database) InsertAlbumIntoDB(album AlbumInfo) (int, error) {
	logrus.WithField("album.Title", album.Title).Debug("Executing InsertAlbumIntoDB")

	query := `INSERT INTO album (title, artist, submitted_on, submitted_by) VALUES(?, ?, ?, ?);`

	logrus.WithField("query", query).Debug("Inserting Query")

	db := self.loadDB()
	defer db.Close()

	res, err := db.Exec(query, album.Title, album.Artist, album.SubmittedOn.Unix(), album.SubmittedBy)
	if err != nil {
		return -1, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return -1, err
	}

	logrus.WithField("id", id).Debug("Album entered into database")
	return int(id), err
}

func (self *Database) loadDB() *sql.DB {
	logrus.WithField("path", self.Path).Debug("Executing loadDB()")

	db, err := sql.Open("sqlite3", self.Path)
	if err != nil {
		panic(err)
	}

	return db
}

func (self *Database) initDB(db *sql.DB) {
	logrus.Debug("Executing initDB()")
	create := `
	CREATE TABLE IF NOT EXISTS album (
		id INTEGER NOT NULL PRIMARY KEY,

		title TEXT,
		artist TEXT,
		year INTEGER,
		submitted_on INTEGER,
		submitted_by TEXT,

		week INTEGER,
		chosen_on INTEGER,
		image TEXT
	);

	CREATE TABLE IF NOT EXISTS selected (
		id INTEGER NOT NULL PRIMARY KEY,

		albumID INTEGER,
		FOREIGN KEY (albumID) REFERENCES album (id)
	);
	`
	_, err := db.Exec(create)
	if err != nil {
		logrus.Fatal(err)
	}
}
