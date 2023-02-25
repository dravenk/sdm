package main

import (
	_ "embed"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
	// "os"
	// "errors"
)

// Config struct
type Config struct {
	Image   string
	Workdir string

	AppsName []string
	Appsdir  string

	Log struct {
		Level int8
	}
	Cmds []string

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

func loadConf() {
	var data []byte
	filePath := "./config.yaml"
	if isNotExist(filePath) {
		data = configfile
	} else {
		data, _ = ioutil.ReadFile("./config.yaml")
	}

	if err := yaml.Unmarshal([]byte(data), &Conf); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
