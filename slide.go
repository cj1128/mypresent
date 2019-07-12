package main

import (
	"io"
	"log"
	"net/http"

	"github.com/cj1128/mypresent/present"
	"github.com/pkg/errors"
)

func handleSlide(w http.ResponseWriter, r *http.Request) {
	err := renderSlide(w, r.URL.Path)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderSlide(w io.Writer, fp string) error {
	doc, err := parseSlide(fp, present.FullMode)

	if err != nil {
		return errors.Wrapf(err, "could not parse slide: %s", fp)
	}

	return doc.Render(w, slideTemplate)
}
