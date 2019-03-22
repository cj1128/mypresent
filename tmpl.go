package main

import (
	"html/template"
	"path/filepath"

	"github.com/fate-lovely/mypresent/present"
)

var (
	indexTemplate *template.Template

	slideTemplate *template.Template
)

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

	tmpl, err = tmpl.New("").Parse(string(actionTmplBytes))
	if err != nil {
		return err
	}

	slideTemplate, err = tmpl.New("").Parse(string(slideTmplBytes))
	if err != nil {
		return err
	}

	dirListTemplateBytes, err := Asset("index.tmpl")
	if err != nil {
		return err
	}

	indexTemplate, err = template.New("").Parse(string(dirListTemplateBytes))
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

		indexTemplate, err = template.ParseFiles(filepath.Join(base, "index.tmpl"))
		return err
	}
}
