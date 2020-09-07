package parser

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Medium/medium-sdk-go"
	"github.com/PuerkitoBio/goquery"
	"github.com/andremueller/md-publisher/config"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type ImageList map[string]*goquery.Selection

func Publish(config config.Config) (*medium.Post, error) {
	client := medium.NewClientWithAccessToken(config.MediumAccessToken)
	u, err := client.GetUser("")
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get medium user")
	}

	log.Infof("User: %+v", u)

	image, err := client.UploadImage(
		medium.UploadOptions{
			FilePath:    "data/figure1.png",
			ContentType: "image/png",
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot upload images to medium")
	}
	log.Infof("Uploaded image: %+v", image)

	p, err := client.CreatePost(medium.CreatePostOptions{
		UserID:        u.ID,
		Title:         "Demo Title 2",
		Content:       fmt.Sprintf("<h1># Demo Title Welcome!\n![caption](%s)\n", image.URL),
		ContentFormat: medium.ContentFormatMarkdown,
		PublishStatus: medium.PublishStatusDraft,
		Tags:          []string{"demo"},
	})
	if err != nil {
		return nil, errors.Wrap(err, "Cannot create medium post")
	}
	log.Infof("Medium post URL: %s\n", p.URL)

	return p, nil
}


func RenderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func RenderDocument(doc *goquery.Document) string {
	return RenderNode(doc.Get(0))
}

func ParseHTML(fileName string) (*goquery.Document, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.WithMessagef(err, "Error while opening file %s", fileName)
	}
	defer file.Close()

	html, err := html.Parse(file)
	if err != nil {
		return nil, errors.WithMessagef(err, "Error while parsing file %s", fileName)
	}

	return goquery.NewDocumentFromNode(html), nil
}

func IsLocalFile(src string) bool {
	return !strings.HasPrefix(src, "http")
}

func FindImages(doc *goquery.Document) ImageList {
	images := make(ImageList)
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && IsLocalFile(src) {
			images[src] = s
		}
	})

	return images
}
