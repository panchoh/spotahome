package model

import (
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

const URL = "https://feeds.spotahome.com/mitula-UK-en.xml"

func Fetch(url string) ([]byte, error) {
	f := path.Base(url)
	cache, err := ioutil.ReadFile(f)
	if err == nil {
		log.Printf("Using cached file: %v", f)
		return cache, nil
	}

	resp, err := http.Get(url)
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
