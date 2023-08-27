package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/alecthomas/kong"
)

var (
	version = "development"
	cli     struct {
		Version    kong.VersionFlag `help:"Print the version and exit" short:"v"`
		Debug      debugFlag        `help:"Enable debug logging."`
		Binaries   string           `help:"Binary files being packaged in separate zip files for lambda." kong:"arg" default:"./bin"`
		OutputPath string           `help:"Where to write the output zip file." kong:"arg" default:"./dist"`
		Template   string           `help:"Template to use for bootstrap file." default:"#!/bin/sh\nexec $LAMBDA_TASK_ROOT/{{ . }}"`
	}
)

type debugFlag bool

func (d debugFlag) BeforeApply(logger *log.Logger) error {
	logger.SetOutput(os.Stderr)
	return nil
}

func main() {
	logger := log.New(io.Discard, "", log.LstdFlags)

	kong.Parse(&cli,
		kong.Description("Packaging tool which builds Lambda deployment archives."),
		kong.Bind(logger),
		kong.Configuration(kong.JSON, ".lambdapack.json"),
		kong.Vars{
			"version": version,
		},
	)

	logger.Println("Starting Lambda packaging for binaries:", cli.Binaries)

	tpl, err := template.New("bootstrap").Parse(cli.Template)
	if err != nil {
		log.Fatal("failed to parse bootstrap template:", err)
	}

	// files, err := os.ReadDir(cli.BinPath)
	files, err := filepath.Glob(cli.Binaries)
	if err != nil {
		log.Fatal("failed to read bin path:", err)
	}

	for _, file := range files {

		fin, err := os.Stat(file)
		if err != nil {
			log.Fatal("failed to stat file:", err)
		}

		if fin.IsDir() {
			continue
		}

		if strings.HasPrefix(fin.Name(), ".") {
			continue
		}

		if strings.HasSuffix(fin.Name(), ".zip") {
			// skip existing zip files
			continue
		}

		logger.Println("Processing file:", fin.Name())

		archivePath := filepath.Join(cli.OutputPath, fin.Name()+".zip")

		logger.Println("Creating archive file:", archivePath)

		archive, err := os.Create(archivePath)
		if err != nil {
			log.Fatal("failed to create archive:", err)
		}

		defer archive.Close()

		zipWriter := zip.NewWriter(archive)

		binf, err := os.Open(file)
		if err != nil {
			log.Fatal("failed to open file:", err)
		}

		defer binf.Close()

		binw, err := zipWriter.Create(fin.Name())
		if _, err := io.Copy(binw, binf); err != nil {
			log.Fatal("failed to write file to archive:", err)
		}

		logger.Println("Writing bootstrap file")

		bootw, err := zipWriter.Create("bootstrap")
		if _, err := io.Copy(binw, binf); err != nil {
			log.Fatal("failed to write file to archive:", err)
		}

		if err = tpl.Execute(bootw, filepath.Base(fin.Name())); err != nil {
			log.Fatal("failed to write bootstrap file:", err)
		}

		logger.Println("Closing archive")
		if err := zipWriter.Close(); err != nil {
			log.Fatal("failed to close archive:", err)
		}
	}
}
