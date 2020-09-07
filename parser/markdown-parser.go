package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func parseMarkdown() {
	source, err := ioutil.ReadFile("data/demo_article.md")
	if err != nil {
		log.Fatal(err)
	}
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			mathjax.MathJax,
			meta.Meta,
		),
	)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert(source, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	metaData := meta.Get(context)
	for k, v := range metaData {
		fmt.Printf("%s -> %+v\n", k, v)
	}
	ioutil.WriteFile("out.html", buf.Bytes(), 0644)
}
