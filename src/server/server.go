package server

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	Database *Database
}

func New(db_path string) *Server {
	s := &Server{}
	logrus.SetLevel(logrus.DebugLevel)
	s.Database = NewDatabase(db_path)
	return s
}

func (s *Server) Run() {
	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/test", HandleTest)
	http.HandleFunc("/search", HandleSearch)
	http.HandleFunc("/submit", s.HandleSubmit)

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	logrus.Info("Starting server")
	http.ListenAndServe(":8080", nil)
}

// GET /
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	logrus.Info("Handling /")
	dummyInfo := Index{
		Album: AlbumInfo{
			Title:  "Test Album",
			Artist: "Test Artist",
			Year:   2023,
			Image:  "static/images/Marquee Moon.jpg",
		},
		Form: Form{},
	}
	t := template.Must(template.ParseGlob("templates/*.html"))
	t.ExecuteTemplate(w, "index", dummyInfo)
}

// GET /test
func HandleTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}
	logrus.Info("Handling /test")
	dummy := Options{
		{Title: "Test1", Artist: "ARtsist 1"},
		{Title: "American Football", Artist: "Americna Footabla"},
	}
	t := template.Must(template.ParseFiles("templates/dummyOpt.html"))

	t.ExecuteTemplate(w, "options", dummy)
}

// GET /search
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}
	logrus.Info("Handling /search")

	title := r.URL.Query().Get("submission")
	logrus.WithField("title", title).Debug("Searching for title")

	// options := queryOptions(title)
	options := dummyOptions()

	t := template.Must(template.ParseFiles("templates/options.html"))
	t.ExecuteTemplate(w, "options", options)
}

// TODO POST /submit
func (s *Server) HandleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	logrus.Info("Handling /submit")

	r.ParseForm()

	formEntry := r.FormValue("submission")
	logrus.WithField("entry", formEntry).Debug("/submit pinged with form entry")

	albumInfo, err := ParseSubmission(formEntry)
	if err != nil {
		http.Error(w,
			"Could not parse submission. Use auto-complete results",
			http.StatusBadRequest)
	}

	fmt.Println(albumInfo)

	albumInfo.SubmittedBy = GetIPHash(r)
	albumInfo.SubmittedOn = time.Now()

	_, err = s.Database.InsertAlbumIntoDB(albumInfo)
	if err != nil {
		logrus.Errorf("Error inserting album: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
}
