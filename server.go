package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"mime"
)

func respond(w http.ResponseWriter, body string) {
	fmt.Fprintf(w, `<!doctype html>
		<head>
		<meta name="viewport" content="initial-scale=1.0,maximum-scale=1.0"/>
		<meta charset="utf-8"/>
		</head>
		<body>%s</body>`, body)
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]

	fi, err := os.Stat(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<h1>404 - Not Fount</h1><p>File: <em>%s</em> not found", path)
		return
	}
	if !fi.IsDir() {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			respond(w, fmt.Sprintf("can't read file %s", path))
			return
		}
		// get content type
		mimeType := mime.TypeByExtension(filepath.Ext(path))
		if mimeType != "" {
			w.Header().Set("content-type", mimeType)
		}
		w.Write(content)
		return
	}
	// directory
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		respond(w, fmt.Sprintf("can't read dir %s", path))
		return
	}
	body := fmt.Sprintf("<h1>Directory %s</h1>", path)
	for _, fi := range fis {
		name := fi.Name()
		file := path + "/" + fi.Name()
		if fi, err := os.Stat(file); err == nil && fi.IsDir() {
			name += "/"
		}
		body += fmt.Sprintf("<li><a href=\"/%s\">%s</a></li>\n", file, name)
	}
	respond(w, body)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
