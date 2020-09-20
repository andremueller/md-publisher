package publisher

import (
	"fmt"
	"path"
	"strings"

	"github.com/Medium/medium-sdk-go"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/andremueller/md-publisher/config"
	"github.com/andremueller/md-publisher/file"
	"github.com/andremueller/md-publisher/parser"
)

// DetectContentType returns the image content type
//  ``image/jpeg``, ``image/png``, ``image/gif``, and ``image/tiff``.
func DetectContentType(fileName string) (string, error) {
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
func PublishMedium(inputFileName string, config config.Config) (*medium.Post, error) {
	if config.MediumAccessToken == "" {
		return nil, fmt.Errorf("Medium access token is empty. Specify it on the command line with --medium-token or put it in the config file")
	}

	doc, err := parser.ParseHTML(inputFileName)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot parse HTML file %s", inputFileName)
	}

	basePath := path.Dir(inputFileName)

	client := medium.NewClientWithAccessToken(config.MediumAccessToken)
	user, err := client.GetUser("")
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get medium user")
	}
	log.Infof("Logged in as user %s", user.Username)

	if !config.NoImages {
		err = processImages(basePath, doc, client)
		if err != nil {
			return nil, errors.Wrapf(err, "Error while processing images for file %s", inputFileName)
		}
	} else {
		log.Info("Skipping images")
	}

	html := parser.RenderDocument(doc)

	post, err := client.CreatePost(medium.CreatePostOptions{
		UserID:        user.ID,
		Title:         parser.GetTitle(doc, "Unknown Title"),
		Content:       html,
		ContentFormat: medium.ContentFormatHTML,
		PublishStatus: medium.PublishStatusDraft,
		Tags:          parser.GetKeywords(doc),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Cannot create medium post")
	}
	log.Infof("Medium post URL: %s\n", post.URL)

	return post, nil
}

func processImages(basePath string, doc *goquery.Document, client *medium.Medium) error {
	images := parser.FindImages(doc)
	for src, locations := range images {
		contentType, err := DetectContentType(src)
		if err != nil {
			return errors.Errorf("Cannot detect content type for image %s", src)
		}
		imageFile := src
		if !path.IsAbs(src) {
			imageFile = path.Join(basePath, imageFile)
		}
		if !file.Exists(imageFile) {
			return fmt.Errorf("Could not find referenced file %s (src = %s)", imageFile, src)
		}
		log.Infof("Uploading image %s [%s]", imageFile, contentType)
		image, err := client.UploadImage(
			medium.UploadOptions{
				FilePath:    imageFile,
				ContentType: contentType,
			},
		)
		if err != nil {
			return errors.Wrapf(err, "Cannot upload image %s to medium", imageFile)
		}
		log.Infof("Mapping image %s in %d location(s) to %s", src, len(locations), image.URL)

		for _, location := range locations {
			location.SetAttr("src", image.URL)
		}
	}
	return nil
}
