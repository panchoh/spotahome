package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
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

func main() {
	xmlValue, err := ioutil.ReadFile("mitula-UK-en.xml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var trovit Trovit

	err = xml.Unmarshal(xmlValue, &trovit)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for index, element := range trovit.Ads {
		fmt.Println(index, element.Title)
	}
}
