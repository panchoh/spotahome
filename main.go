package main

import (
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"log"
	"os"
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

	// Dump an HTML page.
	htmlFile, err := os.Create("index.html")
	if err != nil {
		log.Printf("error: %v", err)
	}
	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Trovit!</title>
	</head>
	<body>
		<table>
			<thead>
				<tr>
					<th>ID</th>
					<th>City</th>
					<th>Title</th>
					<th>Picture</th>
				</tr>
			</thead>
			<tbody>
				{{range $ad := .Ads}}
				<tr>
					<td>{{$ad.Id}}</td>
					<td>{{$ad.City}}</td>
					<td><a href="{{$ad.URL}}">{{$ad.Title}}</a></td>
					<td>
					{{with $pics := .Pictures.Pictures}}
						{{with $pic := (index $pics 0).URL}}
						<a href="{{$pic}}">
							<img src="{{$pic}}"
								alt="{{$ad.Title}}, {{$ad.City}}"
								height="100"
							>
						</a>
						{{end}}
					{{else}}
						Picture not available.
					{{end}}
					</td>
				</tr>
				{{end}}
			</tbody>
		</table>
	</body>
</html>`

	t := template.Must(template.New("trovit").Parse(tpl))

	err = t.Execute(htmlFile, trovit)
	if err != nil {
		log.Printf("error: %v", err)
	}
}
