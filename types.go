package main

import (
	"encoding/xml"
	"os"
	"strings"
)

type Jnlp struct {
	Codebase string   // JNLP's codebase
	Href     string   // JNLP's href
	Title    string   // JNLP's title
	Jars     []string // JNLP's jars
}

type _jnlp struct {
	Codebase    string `xml:"codebase,attr"`
	Href        string `xml:"href,attr"`
	Information struct {
		Title string `xml:"title"`
	} `xml:"information"`
	Resources []struct {
		Jars []struct {
			Href string `xml:"href,attr"`
		} `xml:"jar"`
	} `xml:"resources"`
}

func (j *Jnlp) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var temp _jnlp

	if err := d.DecodeElement(&temp, &start); err != nil {
		return err
	}

	j.Codebase = temp.Codebase
	j.Title = temp.Information.Title
	if len(j.Jars) < 1 {
		j.Jars = make([]string, 0)
	}
	for _, res := range temp.Resources {
		for _, jar := range res.Jars {
			j.Jars = append(j.Jars, jar.Href)
		}
	}
	return nil
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
