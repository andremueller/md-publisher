package publisher

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/andremueller/md-publisher/config"
	"github.com/andremueller/md-publisher/parser"
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

func PublishMedium(inputFileName string, config config.Config) error {
	doc, err := parser.ParseHTML(inputFileName)
	if err != nil {
		log.Fatal("Cannot parse", err)
	}
	images := parser.FindImages(doc)
	for k, v := range images {
		fmt.Printf("- %s\n", k)
		contentType, err := GetContentType(k)
		if err != nil {
			return errors.Errorf("Cannot detect content type for image %s", k)
		}
		v.SetAttr("src", "https://"+k)
		v.SetAttr("content", contentType) // todo replace that
	}
	fmt.Println(parser.RenderDocument(doc))
	return nil
}
