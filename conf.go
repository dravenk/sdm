package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
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
	data, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Printf("file error: %v\n ", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal([]byte(data), &Conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

}
