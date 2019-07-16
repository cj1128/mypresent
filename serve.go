package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/cj1128/mypresent/present"
	"github.com/kataras/golog"
	"github.com/pkg/errors"
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
	content, err := getSlideHTML(r.URL.Path)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(content)
}

func getSlideHTML(path string) ([]byte, error) {
	doc, err := parseSlide(path, present.FullMode)

	if err != nil {
		return nil, errors.Wrap(err, "could not parse slide")
	}

	buf := &bytes.Buffer{}

	err = slideTemplate.Execute(buf, struct {
		*present.Doc
		Template     *template.Template
		NotesEnabled bool
	}{doc, slideTemplate, opts.notesEnabled})

	return buf.Bytes(), err
}

func isSlide(path string) bool {
	return filepath.Ext(path) == ".slide"
}

type indexSlide struct {
	Name string
	Path string
}

type indexData struct {
	Name     string // name of the directory
	Slides   []*indexSlide
	Children []*indexData
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	content, err := getIndexHTML()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(content)
}

func getIndexHTML() ([]byte, error) {
	data, err := scanDir(".")

	if err != nil {
		return nil, errors.Wrap(err, "could not scan dir")
	}

	buf := &bytes.Buffer{}

	if err := indexTemplate.Execute(buf, data); err != nil {
		return nil, errors.Wrap(err, "could not execute template")
	}

	return buf.Bytes(), nil
}

// dir is relative to `contentBase`
// top level Name of indexData is `.`
func scanDir(dir string) (*indexData, error) {
	result := &indexData{
		Name:     filepath.Base(dir),
		Children: make([]*indexData, 0),
	}

	files, err := ioutil.ReadDir(path.Join(opts.contentBase, dir))

	if err != nil {
		return nil, errors.Wrapf(err, "could not read dir: %s", dir)
	}

	for _, f := range files {
		if isSlide(f.Name()) {
			data, err := parseIndexSlide(path.Join(dir, f.Name()))

			if err != nil {
				return nil, err
			}

			result.Slides = append(result.Slides, data)
		}

		if f.IsDir() {
			data, err := scanDir(path.Join(dir, f.Name()))

			if err != nil {
				return nil, err
			}

			result.Children = append(result.Children, data)
		}

		// ignore other files
	}

	return result, nil
}

// fp is relative to contentBase
func parseSlide(fp string, mode present.ParseMode) (*present.Doc, error) {
	f, err := os.Open(path.Join(opts.contentBase, fp))

	if err != nil {
		return nil, errors.Wrapf(err, "could not open file: %s", fp)
	}

	return present.Parse(f, path.Join(opts.contentBase, fp), mode)
}

// fp is relative to contentBase
func parseIndexSlide(fp string) (*indexSlide, error) {
	doc, err := parseSlide(fp, present.TitlesOnly)

	if err != nil {
		return nil, errors.Wrapf(err, "could not parse slide: %s", fp)
	}

	return &indexSlide{
		Name: doc.Title,
		Path: fp,
	}, nil
}
