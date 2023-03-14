package main

import (
	_ "embed"
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

//go:embed projects/drupal/default.config.yaml
var configfile []byte

func initConfigFile() {
	filePath := pathJoin(valuesFile)

	if !isNotExist(filePath) {
		configfile, _ = yaml.Marshal(Conf)
	}
	os.WriteFile(filePath, configfile, os.ModePerm)
}

//go:embed projects/drupal/default.docker-compose.yaml
var dockerComposeFile []byte

func initDockerComposefile() {
	filePath := pathJoin(projectPath, tplCompose)
	if isNotExist(filePath) {
		os.WriteFile(filePath, dockerComposeFile, os.ModePerm)
	}
}

//go:embed projects/drupal/default.settings.php
var settingsFile []byte

func initSettingsfile() {
	filePath := pathJoin(projectPath, tplSettings)
	if isNotExist(filePath) {
		os.WriteFile(filePath, settingsFile, os.ModePerm)
	}
}

func writeFileln(dstFileName string, textSlice []string) {
	if len(textSlice) == 0 {
		return
	}

	f, err := os.OpenFile(dstFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, text := range textSlice {
		if _, err = f.WriteString(text + "\r\n"); err != nil {
			panic(err)
		}
	}
}

func mkDir(dir string, perm os.FileMode) {
	if dir == "" {
		return
	}

	logln("Execute: MkdirAll ", dir)
	if err := os.MkdirAll(dir, perm); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal(err)
		}
	}
}

func availablePort(port int) bool {
	p := strconv.Itoa(port)
	ln, err := net.Listen("tcp", ":"+p)
	if err != nil {
		logln(err)
		return false
	}

	if err = ln.Close(); err != nil {
		logln(err)
		return false
	}

	logln("Available Port: ", port)
	return true
}

func isNotExist(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return true
		}
	}
	// THE FILE EXISTS
	return false
}

func generateDockEnv(appName string) []string {
	port := 0
	textSlice := []string{
		`APP_NAME=` + appName,
		`APP_IMAGE=` + Conf.Image,
		`MARIADB_PASS=` + Conf.MySQL.password,
	}
	for i := Conf.Minport; i < Conf.Maxport; i++ {
		if availablePort(i) {
			port = i
			Conf.Minport = i + 1
			portLn := `APP_PORT=` + strconv.Itoa(i)
			textSlice = append(textSlice, portLn)
			break
		}
	}
	if port == 0 {
		return []string{}
	}
	return textSlice
}

func generateSettings(appName string) []string {
	dbStr := `$databases['default']['default']`
	dbUserStr := dbStr + `['username'] = '` + Conf.MySQL.User + `';`
	dbPassStr := dbStr + `['password'] = '` + Conf.MySQL.password + `';`
	dbHostStr := dbStr + `['host'] = '` + Conf.MySQL.Host + `';`
	portStr := strconv.Itoa(int(Conf.MySQL.Port))
	dbPortStr := dbStr + `['port'] = '` + portStr + `';`
	dbNameStr := dbStr + `['database'] = '` + appName + `';`
	hashSaltStr := `$settings['hash_salt'] = '` + hashSalt() + `';`

	textSlice := []string{
		hashSaltStr,
		dbUserStr,
		dbPassStr,
		dbHostStr,
		dbPortStr,
		dbNameStr,
	}
	return textSlice
}

func pathJoin(elem ...string) string {
	return filepath.Join(elem...)
}
