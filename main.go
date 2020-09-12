package main

import (
	"fmt"
	"os"

	"github.com/andremueller/md-publisher/config"
	"github.com/andremueller/md-publisher/publisher"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	version string // version number (set by the build process see https://stackoverflow.com/questions/11354518/application-auto-build-versioning#11355611)
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
		&cli.BoolFlag{Name: "no-images", Usage: "Does not upload images."},
		&cli.StringFlag{Name: "medium-token", Usage: "Medium.com access token"},
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
		err := config.ReadConfig(context.String("config"), &currentConfig)
		if err != nil {
			return errors.Wrapf(err, "Cannot read configuration file")
		}
		// overwrite with command line arguments
		updateConfig(context, &currentConfig)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func publishCommand(context *cli.Context) error {
	if context.Args().Len() != 1 {
		return fmt.Errorf("publish requires an input file")
	}
	inputFileName := context.Args().Get(0)
	log.Infof("Parsing input file %s", inputFileName)
	return publisher.PublishMedium(inputFileName, currentConfig)
}

func updateConfig(context *cli.Context, config *config.Config) {
	config.NoImages = context.Bool("no-images")
	config.MediumAccessToken = context.String("medium-token")
}
