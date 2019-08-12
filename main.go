package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type Trovit struct {
	XMLName  xml.Name `xml:"trovit" json:"-"`
	Ads      []Ad     `xml:"ad"     json:"ad"`
	SortedBy string   `xml:"-"      json:"-"`
}

type Ad struct {
	XMLName  xml.Name `xml:"ad"       json:"-"`
	Id       int      `xml:"id"       json:"id"`
	URL      string   `xml:"url"      json:"url"`
	Title    string   `xml:"title"    json:"title"`
	City     string   `xml:"city"     json:"city"`
	Pictures Pictures `xml:"pictures" json:"pictures"`
}

type Pictures struct {
	XMLName  xml.Name  `xml:"pictures" json:"-"`
	Pictures []Picture `xml:"picture"  json:"picture,omitempty"`
}

type Picture struct {
	XMLName xml.Name `xml:"picture"       json:"-"`
	URL     string   `xml:"picture_url"   json:"picture_url"`
	Title   string   `xml:"picture_title" json:"picture_title"`
}

var (
	httpAddr = flag.String("http", ":8080", "Listen address")
)

func (t *Trovit) sortBy(s string) *Trovit {
	var st Trovit
	var sorter func(i, j int) bool
	switch s {
	case "id":
		sorter = func(i, j int) bool { return st.Ads[i].Id < st.Ads[j].Id }
	case "city":
		sorter = func(i, j int) bool { return st.Ads[i].City < st.Ads[j].City }
	case "title":
		sorter = func(i, j int) bool { return st.Ads[i].Title < st.Ads[j].Title }
	default:
		// Return immediately, don't sort
		return t
	}
	st.Ads = make([]Ad, len(t.Ads))
	copy(st.Ads, t.Ads)

	sort.Slice(st.Ads, sorter)
	st.SortedBy = s

	return &st
}

func fetchXML() ([]byte, error) {
	xmlValue, err := ioutil.ReadFile("mitula-UK-en.xml")
	if err == nil {
		log.Println("fetchXML: Using cached file")
		return xmlValue, nil
	}

	resp, err := http.Get("https://feeds.spotahome.com/mitula-UK-en.xml")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {
	flag.Parse()

	log.Println("Fetching XML.")
	xmlValue, err := fetchXML()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var trovit Trovit

	err = xml.Unmarshal(xmlValue, &trovit)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	t := template.Must(template.ParseFiles("trovit.html.tmpl"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s := r.FormValue("s")
		err = t.Execute(w, trovit.sortBy(s))
		if err != nil {
			log.Printf("error: %v", err)
		}
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		s := r.FormValue("s")
		j, err := json.Marshal(trovit.sortBy(s))
		if err != nil {
			log.Printf("error marshaling JSON: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(j)
	})

	log.Printf("Starting HTTP server on %s.", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
