package main

import (
	_ "embed"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

// see https://pkg.go.dev/embed
//
//go:embed default.config.yaml
var configfile []byte

func initConfigFile() {
	filePath := "./config.yaml"
	if isNotExist(filePath) {
		ioutil.WriteFile(filePath, configfile, os.ModePerm)
	}
}

//go:embed default.docker-compose.yml
var dockerComposeFile []byte

func initDockerComposefile() {
	filePath := "./docker-compose.yml"
	if isNotExist(filePath) {
		ioutil.WriteFile(filePath, dockerComposeFile, os.ModePerm)
	}
}

//go:embed default.settings.php
var settingsFile []byte

func initSettingsfile() {
	filePath := "./settings.php"
	if isNotExist(filePath) {
		ioutil.WriteFile(filePath, settingsFile, os.ModePerm)
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

	logln("mkdir ", dir)
	if err := os.Mkdir(dir, perm); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal(err)
		}
	}
}

func ScanPort(protocol string, hostname string, port int) bool {
	p := strconv.Itoa(port)
	addr := net.JoinHostPort(hostname, p)
	conn, err := net.DialTimeout(protocol, addr, 2*time.Second)
	if err != nil {
		// logln(err)
		return false
	}
	defer conn.Close()
	return true
}

func portReady(port int) bool {
	if !ScanPort("http", "localhost", port) && !ScanPort("https", "localhost", port) {
		logln("Port available. ", port)
		return true
	}
	return false
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
		if portReady(i) {
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
	// $databases['default']['default']['username'] = 'sqlusername';
	// $databases['default']['default']['password'] = 'sqlpassword';
	// $databases['default']['default']['host'] = 'localhost';
	// $databases['default']['default']['port'] = '3306';
	// $settings['hash_salt'] = '';
	// $databases['default']['default']['database'] = '';
	dbStr := `$databases['default']['default']`
	dbUserStr := dbStr + `['username'] = '` + Conf.MySQL.User + `';`
	// dbPassStr := dbStr + `['password'] = '` + Conf.MySQL.Pass + `';`
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
