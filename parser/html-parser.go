package parser

import (
	"bytes"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

// GoquerySelections is a list of DOM tree selections.
type GoquerySelections []*goquery.Selection

// ImageList is a list of images mapping from a local file name to multiple
// goquery.Selections.
type ImageList map[string]GoquerySelections

// RenderNode renders the html.Node as a string
func RenderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

// RenderDocument renders a goquery.Document as string.
func RenderDocument(doc *goquery.Document) string {
	return RenderNode(doc.Get(0))
}

// ParseHTML parses the given HTML file and returns on success a goquery Document.
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

// IsLocalFile returns true if the given img src refers to a local file.
func IsLocalFile(src string) bool {
	return !strings.HasPrefix(src, "http")
}

// FindImages returns an ImageList of all local images within a DOM tree.
// The ImageList maps the "src" attribute to a list of goquery.Selections
// which can be used for modifying the DOM tree after uploading the image.
func FindImages(doc *goquery.Document) ImageList {
	images := make(ImageList)
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && IsLocalFile(src) {
			src, err := url.QueryUnescape(src)
			if err != nil {
				logrus.Errorf("Cannot decode image path %s", src)
			}
			result, keyExists := images[src]
			if !keyExists {
				result = make(GoquerySelections, 0)
			}
			images[src] = append(result, s)
		}
	})

	return images
}

// SplitList splits a "," separated attribute list into an array
func SplitList(text string) []string {
	v := strings.Split(text, ",")
	result := make([]string, len(v))
	for i, entry := range v {
		result[i] = strings.TrimSpace(entry)
	}
	return result
}

// GetKeywords gets the keywords meta tag from the HTML DOM
func GetKeywords(doc *goquery.Document) []string {
	content := doc.Find("meta[name=keywords]").AttrOr("content", "")
	return SplitList(content)
}

// GetTitle returns the head.title or defaultTitle if not found.
func GetTitle(doc *goquery.Document, defaultTitle string) string {
	t := doc.Find("head title")
	if t.Nodes == nil {
		return defaultTitle
	}
	return t.Text()
	// return t.AttrOr("content", defaultTitle)
}
