package main

import (
	"context"
	"path"
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
		Version    kong.VersionFlag `help:"Print the version and exit" short:"v"`
		Debug      bool             `help:"Enable debug logging."`
		BinPath    string           `help:"Path containing binaries being packaged in zip files for lambda." kong:"arg" default:"./bin"`
		OutputPath string           `help:"Where to write the output zip file." kong:"arg" default:"./dist"`
		Template   string           `help:"Template to use for bootstrap file." default:"#!/bin/sh\nexec $LAMBDA_TASK_ROOT/{{ . }}"`
		BinPattern string           `help:"Glob pattern to match binaries in bin path." default:"*"`
	}
)

func main() {
	kong.Parse(&cli,
		kong.Description("Packaging tool which builds Lambda deployment archives."),
		kong.Vars{
			"version": version,
		},
	)

	logger := zerolog.New(zerolog.NewConsoleWriter()).Level(zerolog.InfoLevel).With().Caller().Timestamp().Logger()
	if cli.Debug {
		logger = logger.Level(zerolog.DebugLevel)
	}

	binPattern := path.Join(cli.BinPath, cli.BinPattern)

	logger.Debug().Str("binPattern", binPattern).Msg("Starting Lambda packaging")

	bootstrapTemplate, err := template.New("bootstrap").Parse(cli.Template)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse bootstrap template")
	}

	files, err := binaries.Glob(binPattern)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read bin path")
	}

	ctx := logger.WithContext(context.Background())

	for _, binaryFile := range files {

		binaryName := filepath.Base(binaryFile)
		archiveFile := filepath.Join(cli.OutputPath, binaryName+".zip")

		err = archiver.PackageFile(ctx, bootstrapTemplate, binaryFile, archiveFile)
		if err != nil {
			log.Fatal().Err(err).Str("binaryFile", binaryFile).Msg("failed to package file")
		}
	}
}
