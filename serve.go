package main

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/cj1128/mypresent/present"
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

	if present.NotesEnabled {
		golog.Info("notes are enabled, press 'N' from the browser to display them.")
	}

	golog.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", opts.host, opts.port), nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// temp, refresh templates everytime
	initTemplates()

	path := r.URL.Path

	if path == "/favicon.ico" {
		http.NotFound(w, r)
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

func isSlide(path string) bool {
	return filepath.Ext(path) == ".slide"
}
