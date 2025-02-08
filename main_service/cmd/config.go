package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	AudioExtractorPort string `envconfig:"AUDIO_EXTRACTOR_PORT" default:"5001"`
	TranscriberPort    string `envconfig:"TRANSCRIBER_PORT" default:"5002"`
	MainServicePort    string `envconfig:"MAIN_SERVICE_PORT" default:"8080"`
	DownloaderPort     string `envconfig:"DOWNLOADER_PORT" default:"5011"`
}

func newConfig() *config {
	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
