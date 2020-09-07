package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/html"

	medium "github.com/Medium/medium-sdk-go"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"

	"github.com/PuerkitoBio/goquery"
)

var (
	version string // version number (set by the build process see https://stackoverflow.com/questions/11354518/application-auto-build-versioning#11355611)
)

// Image ContentType
//  ``image/jpeg``, ``image/png``, ``image/gif``, and ``image/tiff``.

func GetContentType(fileName string) (string, error) {
	f := strings.ToLower(fileName)
	result := ""
	switch {
	case strings.HasSuffix(f, ".png"):
		result = "png"
	case strings.HasSuffix(f, ".jpg") || strings.HasSuffix(f, ".jpeg"):
		result = "jpeg"
	case strings.HasSuffix(f, ".gif"):
		result = "gif"
	case strings.HasSuffix(f, ".tif") || strings.HasSuffix(f, ".tiff"):
		result = "tiff"
	default:
		return "", fmt.Errorf("Unsupported image %s", fileName)
	}
	return "image/" + result, nil
}

type Config struct {
	MediumAccessToken string
}

func ReadConfig(context *cli.Context) Config {
	var configfile = context.String("config")
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}

func publish(config Config) (*medium.Post, error) {
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

func IsLocalFile(src string) bool {
	return !strings.HasPrefix(src, "http")
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

func parseHTML(fileName string) (*goquery.Document, error) {
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

type ImageList map[string]*goquery.Selection

func findImages(doc *goquery.Document) ImageList {
	images := make(ImageList)
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && IsLocalFile(src) {
			images[src] = s
		}
	})

	return images
}

func main() {
	// log.Out = os.Stdout
	log.SetFormatter(&log.TextFormatter{})
	app := cli.NewApp()
	app.Name = "mdpublisher"
	app.Usage = "Publishes an articles to medium.com"
	app.Version = version

	// common flags
	app.Flags = []cli.Flag{
		&cli.IntFlag{Name: "log-level",
			Usage:   "set logging level to (5 = debug, 4 = info, 3 = warn, 2 = error, 1 = fatal",
			Value:   2,
			Aliases: []string{"L"}},
		&cli.StringFlag{Name: "config",
			Usage:   "mdpublisher config file",
			Value:   "~/.config/mdpublisher/mdpublisher.conf",
			Aliases: []string{"c"}},
	}

	publishFlags := []cli.Flag{
		&cli.BoolFlag{Name: "no-images", Usage: "Does not upload images."},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "publish",
			Usage:  "publishs the given article",
			Flags:  publishFlags,
			Action: publishCommand},
	}
	app.Before = func(context *cli.Context) error {
		level := log.Level(context.Int("log-level"))
		log.SetLevel(level)
		log.Errorf("Setting log level to %v", level)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := parseHTML("./data/demo_article.html")
	if err != nil {
		log.Fatal("Cannot parse", err)
	}
	images := findImages(doc)
	for k, v := range images {
		fmt.Printf("- %s\n", k)
		contentType, err := GetContentType(k)
		if err != nil {
			log.Fatalf("Cannot detect content type for image %s", k)
		}
		v.SetAttr("src", "https://"+k)
		v.SetAttr("content", contentType) // todo replace that
	}
	fmt.Println(RenderDocument(doc))
}

func publishCommand(context *cli.Context) error {
	if context.Args().Len() != 1 {
		return fmt.Errorf("publish requires an input file")
	}
	inputFileName := context.Args().Get(0)
	log.Info("Parsing input file %s", inputFileName)

	return nil
}
