package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/panchoh/spotahome/model"
)

var (
	httpAddr = flag.String("http", ":8080", "Listen address")
)

func main() {
	flag.Parse()

	log.Println("Fetching XML.")
	xmlValue, err := model.FetchXML()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var trovit model.Trovit

	err = xml.Unmarshal(xmlValue, &trovit)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	t := template.Must(template.ParseFiles("trovit.html.tmpl"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s := r.FormValue("s")
		err = t.Execute(w, trovit.SortBy(s))
		if err != nil {
			log.Printf("error: %v", err)
		}
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		s := r.FormValue("s")
		j, err := json.Marshal(trovit.SortBy(s))
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