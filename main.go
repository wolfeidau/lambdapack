package main

import (
	"archive/zip"
	"fmt"
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

	logger = log.New(io.Discard, "", log.LstdFlags)
)

type debugFlag bool

func (d debugFlag) BeforeApply(logger *log.Logger) error {
	logger.SetOutput(os.Stderr)
	return nil
}

func main() {
	kong.Parse(&cli,
		kong.Description("Packaging tool which builds Lambda deployment archives."),
		kong.Bind(logger),
		kong.Vars{
			"version": version,
		},
	)

	logger.Println("Starting Lambda packaging for binaries:", cli.Binaries)

	tpl, err := template.New("bootstrap").Parse(cli.Template)
	if err != nil {
		log.Fatal("failed to parse bootstrap template:", err)
	}

	files, err := filepath.Glob(cli.Binaries)
	if err != nil {
		log.Fatal("failed to read bin path:", err)
	}

	for _, file := range files {
		err := packageFile(tpl, file)
		if err != nil {
			log.Fatal("failed to package file:", err)
		}
	}
}

func packageFile(tpl *template.Template, file string) error {
	fin, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	if fin.IsDir() {
		return nil
	}

	if strings.HasPrefix(fin.Name(), ".") {
		return nil
	}

	if strings.HasSuffix(fin.Name(), ".zip") {
		// skip existing zip files
		return nil
	}

	logger.Println("Processing file:", fin.Name())

	archivePath := filepath.Join(cli.OutputPath, fin.Name()+".zip")

	logger.Println("Creating archive file:", archivePath)

	archive, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("failed to create archive: %w", err)
	}

	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	binf, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	defer binf.Close()

	binw, err := zipWriter.Create(fin.Name())
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	if _, err := io.Copy(binw, binf); err != nil {
		return fmt.Errorf("failed to write file to archive: %w", err)
	}

	logger.Println("Writing bootstrap file")

	bootw, err := zipWriter.Create("bootstrap")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	if _, err := io.Copy(binw, binf); err != nil {
		return fmt.Errorf("failed to write file to archive: %w", err)
	}

	if err = tpl.Execute(bootw, filepath.Base(fin.Name())); err != nil {
		return fmt.Errorf("failed to write bootstrap file: %w", err)
	}

	logger.Println("Closing archive")
	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("failed to close archive: %w", err)
	}

	return nil
}
