package publisher

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/andremueller/md-publisher/config"
	"github.com/andremueller/md-publisher/parser"
)

// GetContentType returns the image content type
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

// PublishMedium publishes the given HTML to medium.com
func PublishMedium(inputFileName string, config config.Config) error {
	doc, err := parser.ParseHTML(inputFileName)
	if err != nil {
		log.Fatal("Cannot parse", err)
	}
	images := parser.FindImages(doc)
	for src, locations := range images {
		fmt.Printf("- %s\n", src)
		contentType, err := GetContentType(src)
		if err != nil {
			return errors.Errorf("Cannot detect content type for image %s", src)
		}
		log.Infof("Mapping image %s in %d location(s)", src, len(locations))
		for _, location := range locations {
			location.SetAttr("src", "https://"+src)
			location.SetAttr("content", contentType) // todo replace that
		}
	}
	fmt.Println(parser.RenderDocument(doc))
	return nil
}
