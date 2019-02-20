// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fate-lovely/mypresent/present"
)

func init() {
	http.HandleFunc("/", dirHandler)
}

type dirListData struct {
	Path                          string
	Dirs, Slides, Articles, Other dirEntrySlice
}

type dirEntry struct {
	Name, Path, Title string
}

type dirEntrySlice []dirEntry

func (s dirEntrySlice) Len() int           { return len(s) }
func (s dirEntrySlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s dirEntrySlice) Less(i, j int) bool { return s[i].Name < s[j].Name }

var (
	// dirListTemplate holds the front page template.
	dirListTemplate *template.Template

	slideTemplate *template.Template
)

// dirHandler serves a directory listing for the requested path, rooted at *contentPath.
func dirHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	name := filepath.Join(contentPath, r.URL.Path)

	if isSlide(name) {
		err := renderSlide(w, name)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	if isDir, err := dirList(w, name); err != nil {
		addr, _, e := net.SplitHostPort(r.RemoteAddr)
		if e != nil {
			addr = r.RemoteAddr
		}
		log.Printf("request from %s: %s", addr, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if isDir {
		return
	}

	http.FileServer(http.Dir(contentPath)).ServeHTTP(w, r)
}

func isSlide(path string) bool {
	return filepath.Ext(path) == ".slide"
}

func initTemplatesFromDefaultResource() error {
	var err error
	tmpl := present.Template()

	actionTmplBytes, err := Asset("action.tmpl")
	if err != nil {
		return err
	}

	slideTmplBytes, err := Asset("slide.tmpl")
	if err != nil {
		return err
	}

	slideTemplate, err = tmpl.New("").Parse(string(actionTmplBytes))
	if err != nil {
		return err
	}

	slideTemplate, err = slideTemplate.New("").Parse(string(slideTmplBytes))
	if err != nil {
		return err
	}

	dirListTemplateBytes, err := Asset("dir.tmpl")
	if err != nil {
		return err
	}

	dirListTemplate, err = template.New("").Parse(string(dirListTemplateBytes))
	return err
}

func initTemplates(base string) error {
	if base == "" {
		return initTemplatesFromDefaultResource()
	} else {
		var err error
		tmpl := present.Template()
		slideTemplate, err = tmpl.ParseFiles(
			filepath.Join(base, "action.tmpl"),
			filepath.Join(base, "slide.tmpl"),
		)

		if err != nil {
			return err
		}

		dirListTemplate, err = template.ParseFiles(filepath.Join(base, "dir.tmpl"))
		return err
	}
}

func renderSlide(w io.Writer, docFile string) error {
	// Read the input and build the doc structure.
	doc, err := parse(docFile, present.FullMode)
	if err != nil {
		return err
	}

	// Execute the template.
	return doc.Render(w, slideTemplate)
}

func parse(name string, mode present.ParseMode) (*present.Doc, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return present.Parse(f, name, mode)
}

// dirList scans the given path and writes a directory listing to w.
// It parses the first part of each .slide file it encounters to display the
// presentation title in the listing.
// If the given path is not a directory, it returns (isDir == false, err == nil)
// and writes nothing to w.
func dirList(w io.Writer, name string) (isDir bool, err error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return false, err
	}
	if isDir = fi.IsDir(); !isDir {
		return false, nil
	}
	fis, err := f.Readdir(0)
	if err != nil {
		return false, err
	}
	strippedPath := strings.TrimPrefix(name, filepath.Clean(contentPath))
	strippedPath = strings.TrimPrefix(strippedPath, "/")
	d := &dirListData{Path: strippedPath}
	for _, fi := range fis {
		// skip the golang.org directory
		if name == "." && fi.Name() == "golang.org" {
			continue
		}
		e := dirEntry{
			Name: fi.Name(),
			Path: filepath.ToSlash(filepath.Join(strippedPath, fi.Name())),
		}
		if fi.IsDir() && showDir(e.Name) {
			d.Dirs = append(d.Dirs, e)
			continue
		}
		if isSlide(e.Name) {
			fn := filepath.ToSlash(filepath.Join(name, fi.Name()))
			if p, err := parse(fn, present.TitlesOnly); err != nil {
				log.Printf("parse(%q, present.TitlesOnly): %v", fn, err)
			} else {
				e.Title = p.Title
			}
			switch filepath.Ext(e.Path) {
			case ".article":
				d.Articles = append(d.Articles, e)
			case ".slide":
				d.Slides = append(d.Slides, e)
			}
		} else if showFile(e.Name) {
			d.Other = append(d.Other, e)
		}
	}
	if d.Path == "." {
		d.Path = ""
	}
	sort.Sort(d.Dirs)
	sort.Sort(d.Slides)
	sort.Sort(d.Articles)
	sort.Sort(d.Other)
	return true, dirListTemplate.Execute(w, d)
}

// showFile reports whether the given file should be displayed in the list.
func showFile(n string) bool {
	switch filepath.Ext(n) {
	case ".pdf":
	case ".html":
	case ".go":
	default:
		return isSlide(n)
	}
	return true
}

// showDir reports whether the given directory should be displayed in the list.
func showDir(n string) bool {
	if len(n) > 0 && (n[0] == '.' || n[0] == '_') || n == "present" {
		return false
	}
	return true
}
