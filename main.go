package main

import (
	"encoding/xml"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Trovit struct {
	XMLName xml.Name `xml:"trovit"`
	Ads     []Ad     `xml:"ad"`
}

type Ad struct {
	XMLName  xml.Name `xml:"ad"`
	Id       int      `xml:"id"`
	URL      string   `xml:"url"`
	Title    string   `xml:"title"`
	City     string   `xml:"city"`
	Pictures Pictures `xml:"pictures"`
}

type Pictures struct {
	XMLName  xml.Name  `xml:"pictures"`
	Pictures []Picture `xml:"picture"`
}

type Picture struct {
	XMLName xml.Name `xml:"picture"`
	URL     string   `xml:"picture_url"`
	Title   string   `xml:"picture_title"`
}

var (
	httpAddr = flag.String("http", ":8080", "Listen address")
)

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
		err = t.Execute(w, trovit)
		if err != nil {
			log.Printf("error: %v", err)
		}
	})

	log.Printf("Starting HTTP server on %s.", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
