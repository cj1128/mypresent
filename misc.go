package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func initTemplate(path string, parent *template.Template) (*template.Template, error) {
	buf, err := getAsset(path)

	if err != nil {
		return nil, errors.Wrap(err, "could not get asset")
	}

	return parent.New("").Parse(string(buf))
}

// get asset from user specified resource dir
// if not present, get it from bundled assets
func getAsset(path string) ([]byte, error) {
	if opts.resourcePath != "" && fileExists(filepath.Join(opts.resourcePath, path)) {
		return ioutil.ReadFile(filepath.Join(opts.resourcePath, path))
	}

	return assetBox.Find(path)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
