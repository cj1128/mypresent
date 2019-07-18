package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kataras/golog"
	"github.com/pkg/errors"
)

func buildContent() {
	// helper functions
	copy := func(path string) {
		input, err := ioutil.ReadFile(filepath.Join(opts.contentBase, path))
		if err != nil {
			golog.Fatal(err)
		}

		err = ioutil.WriteFile(filepath.Join(opts.output, path), input, 0644)
		if err != nil {
			golog.Fatal(err)
		}
	}

	// path is relative to `output`
	mkdir := func(path string) {
		if err := os.MkdirAll(filepath.Join(opts.output, path), 0755); err != nil {
			golog.Fatalf("could not create dir: %v", err)
		}
	}

	// path is relative to `opts.output`
	// fatal if error
	write := func(path string, content []byte) {
		p := filepath.Join(opts.output, path)
		if err := ioutil.WriteFile(p, content, 0644); err != nil {
			golog.Fatal(err)
		}
	}

	// change `.slide` -> `.html`
	modifyPath := func(path string) string {
		return path[:len(path)-len(".slide")] + ".html"
	}

	// create dir
	mkdir(".")

	if err := filepath.Walk(opts.contentBase, func(p string, info os.FileInfo, err error) error {
		path, _ := filepath.Rel(opts.contentBase, p)

		// skip the top level
		if path == "." {
			return nil
		}

		// create dir
		if info.IsDir() {
			mkdir(path)
			return nil
		}

		// generate htmls for slide
		if isSlide(path) {
			content, err := getSlideHTML(path)

			if err != nil {
				return errors.Wrapf(err, "could not render slide: %s", path)
			}

			write(modifyPath(path), content)

			return nil
		}

		// for other files, just copy
		copy(path)

		return nil
	}); err != nil {
		golog.Fatal(err)
	}

	// copy static resources
	mkdir("static")
	for _, path := range assetBox.List() {
		content, err := getAsset(path)
		if err != nil {
			golog.Fatal(err)
		}

		p := filepath.Join("static", path)
		mkdir(filepath.Dir(p))
		write(p, content)
	}

	// generate index.html
	data, err := scanDir(".")
	if err != nil {
		golog.Fatal(err)
	}

	allSlides := getAllSlides(data)

	for _, slide := range allSlides {
		slide.Path = modifyPath(slide.Path)
	}

	buf := &bytes.Buffer{}
	if err := indexTemplate.Execute(buf, struct {
		Index *indexData
		All   []*slideData
	}{data, allSlides}); err != nil {
		golog.Fatal(err)
	}

	write("index.html", buf.Bytes())
}
