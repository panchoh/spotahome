package model

import (
	"encoding/xml"
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

func (t *Trovit) SortBy(s string) *Trovit {
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

func FetchXML() ([]byte, error) {
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
