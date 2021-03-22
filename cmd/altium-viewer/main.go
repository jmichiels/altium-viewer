package main

import (
	"fmt"
	"github.com/jmichiels/altium-viewer/pkg/altium"
	"io/ioutil"
	stdlog "log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s my-altium-project.zip\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	// Discard logs from lorca.
	stdlog.SetOutput(ioutil.Discard)
	log := stdlog.New(os.Stdout, "", stdlog.Ltime)
	projectFilePath := os.Args[1]
	projectFile, err := os.Open(projectFilePath)
	if err != nil {
		log.Fatal(err)
	}
	projectFileName := filepath.Base(projectFilePath)
	log.Printf("upload \"%s\" to altium.com\n", projectFileName)
	id, err := altium.UploadProject(projectFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := projectFile.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("upload done, open viewer\n")
	if err := altium.OpenProject(id, projectFileName); err != nil {
		log.Fatal(err)
	}
}
