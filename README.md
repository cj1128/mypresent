# Mypresent

A tool forked from `golang/tools/present`, made following modifications:

- Support syntax highlight
- Can generate static html files
- Indented code can specify language
- Remove all functions except for slide

```bash
$ mypresent --help
usage: mypresent [<flags>] <command> [<args> ...]

Flags:
  -h, --help               Show context-sensitive help (also try --help-long and
                           --help-man).
  -r, --resource=RESOURCE  static resource path, if not provided, use builtin
                           resource
  -c, --content="."        presentation content path

Commands:
  help [<command>...]
    Show help.

  serve* [<flags>]
    Start the server

  build [<flags>]
    Generate output
```

## Slide Format

title
[subtitle]
[time](format: "15:04 2 Jan 2006" or "2 Jan 2006")
[cover image](format: .cover [url])
<blank>
[misc info]
[sections]

## Static Resource

- index.css
- note.js
- slide.js
- slide.css
- favicon.ico

tmpl
  - index.tmpl
  - slide.tmpl

hljs
  - hljs.js
  - hljs.css

## Syntax Highlight

Use [highlight.js](https://highlightjs.org) to do syntax highlight.

See all supported languages in [css-class-reference](https://highlightjs.readthedocs.io/en/latest/css-classes-reference.html).
