package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	port   = flag.String("port", ":9090", "port to use")
	tpl    *template.Template
	dialer *gomail.Dialer
)

func main() {
	flag.Parse()
	files, err := readFiles("./account/templates/")
	if err != nil {
		logrus.Fatalln(err)
	}

	_ = files

	dialer = gomail.NewDialer("smtp.gmail.com", 587, "antibug.ke@gmail.com", "@@antibug2020")

	tpl = template.Must(template.ParseFiles(files...))

	http.HandleFunc("/email", sendEmail)
	http.HandleFunc("/web", render)

	logrus.Infof("Server started on port %s", *port)
	http.ListenAndServe(":"+strings.TrimPrefix(*port, ":"), nil)
}

func sendEmail(w http.ResponseWriter, r *http.Request) {
	view := r.URL.Query().Get("view")
	if view == "" {
		view = "test"
	}
	buffer := bytes.NewBuffer(make([]byte, 0, 64))
	err := tpl.ExecuteTemplate(buffer, "base", &renderData{
		FirstName:      "Gideon",
		LastName:       "Kamau",
		AccountID:      uuid.New().String(),
		Link:           "google.com",
		Token:          fmt.Sprint(rand.Intn(5)),
		Label:          "user",
		WebsiteURL:     "https://twitter.com",
		AppName:        "umrs Network",
		AppDescription: randomdata.Paragraph(),
		PrimaryColor:   "#42A5F5",
		SecondaryColor: "#4CAF50",
		TemplateName:   view,
		Reason:         randomdata.Paragraph(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send Email
	m := gomail.NewMessage()

	m.SetHeader("From", dialer.Username)
	m.SetHeader("To", "gideonhacer@gmail.com")
	m.SetHeader("Subject", "Testing")
	m.SetBody("text/html", buffer.String())

	err = dialer.DialAndSend(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func render(w http.ResponseWriter, r *http.Request) {
	view := r.URL.Query().Get("view")
	if view == "" {
		view = "create"
	}
	err := tpl.ExecuteTemplate(w, "base", &renderData{
		FirstName:      "Gideon",
		LastName:       "Kamau",
		AccountID:      uuid.New().String(),
		Link:           "google.com",
		Token:          fmt.Sprint(rand.Intn(5)),
		Label:          "user",
		WebsiteURL:     "https://twitter.com",
		AppName:        "umrs Network",
		AppDescription: randomdata.Paragraph(),
		PrimaryColor:   "#42A5F5",
		SecondaryColor: "#4CAF50",
		TemplateName:   view,
		Reason:         randomdata.Paragraph(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func readFiles(dir string) ([]string, error) {
	var allFiles []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, filepath.Join(dir, filename))
		}
	}
	return allFiles, nil
}

type renderData struct {
	FirstName      string
	LastName       string
	AccountID      string
	Link           string
	Token          string
	Label          string
	WebsiteURL     string
	AppName        string
	AppDescription string
	PrimaryColor   string
	SecondaryColor string
	TemplateName   string
	Reason         string
}
