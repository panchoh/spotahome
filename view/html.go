package view

import (
	"html/template"
	"log"
)

var Trovit *template.Template

const (
	tmplfile = "trovit.html.tmpl"
	tmpl     = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Trovit!</title>
		<style type="text/css">
			th, td {padding: 15px;}
			th {background-color: #4CAF50;}
			tr:nth-child(even) {background-color: #f2f2f2;}
			.wrap {
				width: 800px;
				margin: 0 auto;
			}
		</style>
	</head>
	<body>
		<div class="wrap">
			<table>
				<thead>
					<tr>
						<th><a href="?s=id">ID</a></th>
						<th><a href="?s=city">City</a></th>
						<th><a href="?s=title">Title</a></th>
						<th>Picture</th>
						<th><a href="/json?s={{.SortedBy}}">JSON</a></th>
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
		</div>
	</body>
</html>`
)

func init() {
	var err error
	Trovit, err = template.ParseFiles(tmplfile)
	if err == nil {
		log.Printf("Using external template: ‘%s’", tmplfile)
		return
	}
	Trovit = template.Must(template.New(tmplfile).Parse(tmpl))
	log.Printf("Using internal template: ‘%s’", tmplfile)
}
