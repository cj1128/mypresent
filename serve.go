package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/cj1128/mypresent/present"
	"github.com/k0kubun/pp"
	"github.com/kataras/golog"
)

func serveContent() {
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/static/")
		content, err := assetBox.Find(name)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.ServeContent(w, r, name, time.Now(), bytes.NewReader(content))
	})

	http.HandleFunc("/", mainHandler)

	golog.Infof("server started, port: %d, host: %s", opts.port, opts.host)

	if opts.notesEnabled {
		golog.Info("notes are enabled, press 'N' from the browser to display them.")
	}

	golog.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", opts.host, opts.port), nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// temp, refresh templates everytime
	initTemplates()

	path := r.URL.Path

	if path == "/favicon.ico" {
		buf, _ := getAsset("static/favicon.ico")
		http.ServeContent(w, r, "favicon.ico", time.Now(), bytes.NewReader(buf))
		return
	}

	if path == "/" || path == "/index.html" {
		handleIndex(w, r)
		return
	}

	if isSlide(path) {
		handleSlide(w, r)
		return
	}

	http.NotFound(w, r)
}

func handleSlide(w http.ResponseWriter, r *http.Request) {
	doc, err := parseSlide(r.URL.Path, present.FullMode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pp.Println(doc)

	slideTemplate.Execute(w, struct {
		*present.Doc
		Template     *template.Template
		NotesEnabled bool
	}{doc, slideTemplate, opts.notesEnabled})
}

func isSlide(path string) bool {
	return filepath.Ext(path) == ".slide"
}
