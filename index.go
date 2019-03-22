package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/fate-lovely/mypresent/present"
)

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
	data, err := scanDir(".")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	indexTemplate.Execute(w, data)
}

// dir is relative to `contentBase`
func scanDir(dir string) (*indexData, error) {
	result := &indexData{
		Name:     filepath.Base(dir),
		Children: make([]*indexData, 0),
	}

	files, err := ioutil.ReadDir(path.Join(contentBase, dir))

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
func parseIndexSlide(fp string) (*indexSlide, error) {
	f, err := os.Open(path.Join(contentBase, fp))

	if err != nil {
		return nil, errors.Wrapf(err, "could not open file: %s", fp)
	}

	doc, err := present.Parse(f, fp, present.TitlesOnly)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse slide: %s", fp)
	}

	return &indexSlide{
		Name: doc.Title,
		Path: fp,
	}, nil
}
