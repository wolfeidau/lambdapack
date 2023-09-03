package main

import (
	"context"
	"path/filepath"
	"text/template"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/lambdapack/internal/archiver"
	"github.com/wolfeidau/lambdapack/internal/binaries"
)

var (
	version = "development"
	cli     struct {
		Version  kong.VersionFlag `help:"Print the version and exit" short:"v"`
		Debug    bool             `help:"Enable debug logging."`
		Binaries []string         `help:"Binaries being packaged in zip files for lambda." kong:"arg"`
		Output   string           `help:"Directory path to write the output zip archives." default:"./dist"`
		Template string           `help:"Template to use for bootstrap file." default:"#!/bin/sh\nexec $LAMBDA_TASK_ROOT/{{ . }}"`
	}
)

func main() {
	kong.Parse(&cli,
		kong.Description("Packaging tool which builds Lambda deployment archives from a list of binaries."),
		kong.Vars{
			"version": version,
		},
	)

	logger := zerolog.New(zerolog.NewConsoleWriter()).Level(zerolog.InfoLevel).With().Caller().Timestamp().Logger()
	if cli.Debug {
		logger = logger.Level(zerolog.DebugLevel)
	}

	logger.Debug().Msg("Starting Lambda packaging")

	bootstrapTemplate, err := template.New("bootstrap").Parse(cli.Template)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse bootstrap template")
	}

	ctx := logger.WithContext(context.Background())

	filtered, err := binaries.Filter(cli.Binaries...)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to filter binaries")
	}

	for _, binaryFile := range filtered {

		binaryName := filepath.Base(binaryFile)
		archiveFile := filepath.Join(cli.Output, binaryName+".zip")

		err = archiver.PackageFile(ctx, bootstrapTemplate, binaryFile, archiveFile)
		if err != nil {
			log.Fatal().Err(err).Str("binaryFile", binaryFile).Msg("failed to package file")
		}
	}
}
