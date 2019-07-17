package main

import (
	"html/template"

	"github.com/cj1128/mypresent/present"
	"github.com/gobuffalo/packr"
	"github.com/kataras/golog"
	"gopkg.in/alecthomas/kingpin.v2"
)

var assetBox = packr.NewBox("./static")

var (
	opts struct {
		host         string
		port         int
		resourcePath string
		contentBase  string
		output       string
		notesEnabled bool
	}

	indexTemplate *template.Template

	slideTemplate *template.Template
)

func parseFlags() string {
	// common flags
	kingpin.Flag("resource", "static resource path, if not provided, use builtin resource").
		Short('r').
		StringVar(&opts.resourcePath)

	kingpin.Flag("content", "presentation content path").
		Short('c').
		Default(".").
		StringVar(&opts.contentBase)

	// serve flags
	serve := kingpin.Command("serve", "Start the server").Default()
	serve.Flag("host", "server host").
		Default("127.0.0.1").
		StringVar(&opts.host)

	serve.Flag("port", "server port").
		Short('p').
		Default("3999").
		IntVar(&opts.port)

	serve.Flag("notes", "enable presenter notes (press 'N' to display").
		Short('n').
		Default("false").
		BoolVar(&opts.notesEnabled)

	// build flags
	build := kingpin.Command("build", "Generate output")
	build.Flag("output", "output path").
		Short('o').
		Default("dist").
		StringVar(&opts.output)

	kingpin.HelpFlag.Short('h')

	return kingpin.Parse()
}

func initTemplates() {
	var err error
	parent := present.Template()

	slideTemplate, err = initTemplate("tmpl/slide.tmpl", parent)
	if err != nil {
		golog.Fatal(err)
	}

	indexTemplate, err = initTemplate("tmpl/index.tmpl", parent)
	if err != nil {
		golog.Fatal(err)
	}
}

func main() {
	cmd := parseFlags()

	initTemplates()

	switch cmd {
	case "serve":
		serveContent()

	case "build":
		buildContent()
	}
}
