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

var opts struct {
	host         string
	port         int
	resourcePath string
	contentBase  string
}

func parseFlags() {
	kingpin.Flag("host", "server host").
		Default("127.0.0.1").
		StringVar(&opts.host)

	kingpin.Flag("port", "server port").
		Short('p').
		Default("3999").
		IntVar(&opts.port)

	kingpin.Flag("resource", "static resource path, if not provided, use builtin resource").
		Short('r').
		StringVar(&opts.resourcePath)

	kingpin.Flag("content", "presentation content path").
		Short('c').
		Default(".").
		StringVar(&opts.contentBase)

	kingpin.Flag("notes", "enable presenter notes (press 'N' to display").
		Default("false").
		BoolVar(&present.NotesEnabled)

	kingpin.HelpFlag.Short('h')

	kingpin.Parse()
}

func main() {
	parseFlags()

	if err := initTemplates(opts.resourcePath); err != nil {
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

	http.HandleFunc("/", mainHandler)

	log.Printf("server started, host: %s, port: %d", opts.host, opts.port)

	if present.NotesEnabled {
		log.Println("notes are enabled, press 'N' from the browser to display them.")
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", opts.host, opts.port), nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// temp, refresh templates everytime
	initTemplates(opts.resourcePath)

	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		handleIndex(w, r)
		return
	}

}
