package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"

	medium "github.com/Medium/medium-sdk-go"
	log "github.com/sirupsen/logrus"
	// "gopkg.in/yaml.v2"
)

// Image ContentType
//  ``image/jpeg``, ``image/png``, ``image/gif``, and ``image/tiff``.

func publish() {
	client := medium.NewClientWithAccessToken("2b5c387ba9160e0f9d488e1606955e8d3466279a5947a94842a52e730c76e98be")
	u, err := client.GetUser("")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User: %+v\n", u)

	/*	data, err := ioutil.ReadFile("data/figure1.png")
		if err != nil {
			log.Fatal(err)
		}*/
	image, err := client.UploadImage(
		medium.UploadOptions{
			FilePath:    "data/figure1.png",
			ContentType: "image/png",
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Image: %+v", image)

	// client.UploadImage(medium.Sco)
	// Create a draft post.
	p, err := client.CreatePost(medium.CreatePostOptions{
		UserID:        u.ID,
		Title:         "Demo Title",
		Content:       "# Demo Title\nWelcome!",
		ContentFormat: medium.ContentFormatMarkdown,
		PublishStatus: medium.PublishStatusDraft,
		Tags:          []string{"demo"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("URL: %s\n", p.URL)
}

func NewImageExtractor() ast.NodeVisitorFunc {
	return func(node ast.Node, entering bool) ast.WalkStatus {
		fmt.Printf("type: %T", node)
		return ast.GoToNext
	}
}

func main() {
	data, err := ioutil.ReadFile("data/demo_article.md")
	if err != nil {
		log.Fatal(err)
	}
	extensions := parser.CommonExtensions
	parser := parser.NewWithExtensions(extensions)
	doc := markdown.Parse(data, parser)
	ast.PrintWithPrefix(os.Stdout, doc, " ")
	ast.Walk(doc, NewImageExtractor())
}
