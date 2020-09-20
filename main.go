package main

import (
	"fmt"
	"os"

	"github.com/andremueller/md-publisher/config"
	"github.com/andremueller/md-publisher/publisher"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	version string // version number (set by the build process see https://stackoverflow.com/questions/11354518/application-auto-build-versioning#11355611)
)

const (
	configNoImages          = "no-images"
	configMediumAccessToken = "medium-token"
)

var currentConfig config.Config

func main() {
	log.SetFormatter(&log.TextFormatter{})
	app := cli.NewApp()
	app.Name = "md-publisher"
	app.Usage = "Publishes an articles to medium.com"
	app.Version = version

	// common flags
	app.Flags = []cli.Flag{
		&cli.IntFlag{Name: "log-level",
			Usage:   "set logging level to (5 = debug, 4 = info, 3 = warn, 2 = error, 1 = fatal",
			Value:   5,
			Aliases: []string{"L"}},
		&cli.StringFlag{Name: "config",
			Usage:   "md-publisher config file",
			Value:   config.GetDefaultConfigFile(),
			Aliases: []string{"c"}},
	}

	publishFlags := []cli.Flag{
		&cli.BoolFlag{Name: configNoImages,
			Usage: "Do not upload images",
			Value: false},
		&cli.StringFlag{Name: configMediumAccessToken,
			Usage: "Medium.com access token - alternative to the configuration file",
			Value: ""},
	}

	app.Commands = []*cli.Command{
		{
			Name:   "publish",
			Usage:  "publish the given article",
			Flags:  publishFlags,
			Action: publishCommand},
	}
	app.Before = func(context *cli.Context) error {
		level := log.Level(context.Int("log-level"))
		log.SetLevel(level)
		configFile := context.String("config")
		err := config.ReadConfig(configFile, &currentConfig)
		if err != nil {
			log.Errorf("Cannot read configuration file %s", configFile)
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func publishCommand(context *cli.Context) error {
	// overwrite configuration parameters with command line arguments
	updateConfig(context, &currentConfig)

	if context.Args().Len() != 1 {
		return fmt.Errorf("publish requires an input file")
	}
	inputFileName := context.Args().Get(0)
	log.Infof("Parsing input file %s", inputFileName)
	_, err := publisher.PublishMedium(inputFileName, currentConfig)
	return err
}

func updateConfig(context *cli.Context, config *config.Config) {
	if context.IsSet(configNoImages) {
		config.NoImages = context.Bool(configNoImages)
	}
	if context.IsSet(configMediumAccessToken) {
		config.MediumAccessToken = context.String(configMediumAccessToken)
	}
}
