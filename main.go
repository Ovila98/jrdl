package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultDownloadDir     string = "downloads"
	FailedDownloadFlag     string = "--failed-download-exit"
	FailedFileCreationFlag string = "--failed-file-creation-exit"
	FailedFileWriteFlag    string = "--failed-file-write-exit"
)

const usage string = "JNLP resource downloader\n\n" +

	"This utility downloads every JAR files listed in the provided JNLP file.\n\n" +

	"usage: jrdl [OPTIONS] <jnlp-file> [download-dir]\n\n" +

	"ARGUMENTS:\n" +
	"<jnlp-file>			Path to the JNLP file to download (required)\n" +
	"[download-dir]		Directory to download the JAR files to (optional)\n\n" +

	"OPTIONS:\n" +
	"--failed-download-exit			Exit with non-zero code if a download fails\n" +
	"--failed-file-creation-exit		Exit with non-zero code if a file cannot be created\n" +
	"--failed-file-write-exit		Exit with non-zero code if a file cannot be written\n\n" +

	"EXAMPLE:\n" +
	"jrdl jnlp-file.jnlp JnlpOutDir\n\n" +

	"REMARKS:\n" +

	"<jnlp-file> and [download-dir] are positional arguments and must be specified in this order regardless of their position among the possibly specified options.\n" +

	"Any command line argument starting with a dash (-) is considered an option and will be ignored if not supported."

func resolveUrl(base, href string) string {
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(href, "/")
}

func downloadJar(jarUrl string) ([]byte, error) {
	r, err := http.Get(jarUrl)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return io.ReadAll(r.Body)
}

type Flags struct {
	InputFile                string
	DownloadDir              string
	IgnoreFailedDownload     bool
	IgnoreFailedFileCreation bool
	IgnoreFailedFileWrite    bool
}

func (f *Flags) Parse() {
	// Defaults
	f.IgnoreFailedDownload = true
	f.IgnoreFailedFileCreation = true
	f.IgnoreFailedFileWrite = true

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			switch arg {
			case FailedDownloadFlag:
				f.IgnoreFailedDownload = false
			case FailedFileCreationFlag:
				f.IgnoreFailedFileCreation = false
			case FailedFileWriteFlag:
				f.IgnoreFailedFileWrite = false
			}
		} else if f.InputFile == "" {
			f.InputFile = arg
		} else if f.DownloadDir == "" {
			f.DownloadDir = arg
		}
	}

	if f.DownloadDir == "" {
		f.DownloadDir = defaultDownloadDir
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(0)
	}

	var f Flags
	f.Parse()

	b, err := os.ReadFile(f.InputFile)
	if err != nil {
		log.Fatalf("cannot open specified jnlp file '%s':\n%s", f.InputFile, err.Error())
	}

	var j Jnlp

	err = xml.Unmarshal(b, &j)
	if err != nil {
		log.Fatalf("cannot unmarshal jnlp file '%s':\n=> %s", f.InputFile, err.Error())
	}

	if len(j.Jars) < 1 {
		log.Printf("no jars found in jnlp file '%s'", f.InputFile)
		os.Exit(0)
	}

	f.DownloadDir = filepath.Join(f.DownloadDir, j.Title)
	_, err = os.Stat(f.DownloadDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(f.DownloadDir, os.ModePerm)
			if err != nil {
				log.Fatalf("cannot create download directory '%s':\n=> %s", f.DownloadDir, err.Error())
			}
		} else {
			fallbackPath := filepath.Join(defaultDownloadDir, j.Title)
			log.Printf("provided download directory '%s' cannot be set, using default '%s'\n", f.DownloadDir, fallbackPath)
			f.DownloadDir = fallbackPath
			_, err = os.Stat(f.DownloadDir)
			if err != nil {
				err = os.MkdirAll(f.DownloadDir, os.ModePerm)
				if err != nil {
					log.Fatalf("cannot create download directory '%s':\n=> %s", f.DownloadDir, err.Error())
				}
			}
		}
	}

	// TODO: Use concurrency based on number of Jars
	for _, jar := range j.Jars {
		jarUrl := resolveUrl(j.Codebase, jar)

		// Download the jar
		jarBytes, err := downloadJar(jarUrl)
		if err != nil {
			log.Printf("cannot download jar '%s' (%s): \n=> %s\n", jar, jarUrl, err.Error())
			if !f.IgnoreFailedDownload {
				log.Fatalln("=> exiting program...")
			}
			log.Printf("=> to exit immediately, use '%s' flag\n", FailedDownloadFlag)
			continue
		}

		jarName := filepath.Base(jar)

		// Create the jar file
		file, err := os.Create(filepath.Join(f.DownloadDir, jarName))
		if err != nil {
			log.Printf("cannot create jar file '%s':\n=> %s\n", jarName, err.Error())
			if !f.IgnoreFailedFileCreation {
				log.Fatalln("=> exiting program...")
			}
			log.Printf("=> to exit immediately, use '%s' flag\n", FailedFileCreationFlag)
			continue
		}
		defer file.Close()

		// Write the jar file
		_, err = file.Write(jarBytes)
		if err != nil {
			log.Printf("cannot write jar file '%s':\n=> %s\n", jarName, err.Error())
			if !f.IgnoreFailedFileWrite {
				log.Fatalln("=> exiting program...")
			}
			log.Printf("=> to exit immediately, use '%s' flag\n", FailedFileWriteFlag)
			continue
		}
	}

	log.Printf("JAR files downloaded successfully to '%s'\n", f.DownloadDir)
}
