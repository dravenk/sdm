package main

import (
	_ "embed"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config struct
type Config struct {
	Image       string
	ContainerId string
	Workdir     string

	AppsName []string
	Appsdir  string

	Log struct {
		Level int8
	}

	MySQL struct {
		Host     string
		Port     int16
		User     string
		Pass     string
		password string
	}
	Minport int
	Maxport int
}

// Conf convert config to Conf variable
var Conf Config

var valuesFile = "config.yaml"

func loadConf() {
	var data []byte
	filePath := pathJoin(valuesFile)
	if isNotExist(filePath) {
		data = configfile
	} else {
		data, _ = os.ReadFile(filePath)
	}

	if err := yaml.Unmarshal([]byte(data), &Conf); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
