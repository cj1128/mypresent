package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fate-lovely/mypresent/present"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	host         string
	port         int
	resourcePath string
	contentPath  string
)

func parseFlags() {
	kingpin.Flag("host", "server host").
		Short('h').
		Default("127.0.0.1").
		StringVar(&host)

	kingpin.Flag("port", "server port").
		Short('p').
		Default("3999").
		IntVar(&port)

	kingpin.Flag("resource", "static resource path, if not provided, use builtin resource").
		Short('r').
		StringVar(&resourcePath)

	kingpin.Flag("content", "presentation content path").
		Short('c').
		Default(".").
		StringVar(&contentPath)

	kingpin.Flag("play", "enable playground").
		Default("true").
		BoolVar(&present.PlayEnabled)

	kingpin.Flag("notes", "enable presenter notes (press 'N' to display").
		Default("false").
		BoolVar(&present.NotesEnabled)

	kingpin.Parse()
}

func main() {
	parseFlags()

	if err := initTemplates(resourcePath); err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/static/")
		content, err := Asset(name)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.ServeContent(w, r, name, time.Now(), bytes.NewReader(content))
	})

	log.Printf("server started, host: %s, port: %d", host, port)

	if present.NotesEnabled {
		log.Println("notes are enabled, press 'N' from the browser to display them.")
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
}
