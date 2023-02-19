package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

var (
	tplSettings = "settings.php"
	tplCompose  = "docker-compose.yml"
)

func init() {
	initConfigFile()
	initDockerComposefile()
	initSettingsfile()

	loadConf()
	logln("----- start ------")
}

func main() {
	mkDir(Conf.Appsdir)

	appsName := []string{"dp"}
	if len(Conf.AppsName) > 0 {
		appsName = Conf.AppsName
	}

	// Create container: docker create drupal:latest
	logln("Execute: docker create", Conf.Image)
	cidcmd := exec.Command("docker", "create", Conf.Image)
	cid, err := cidcmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	containerid := strings.TrimSuffix(string(cid), "\n")
	imgworkdir := containerid + `:` + Conf.Workdir

	for i := 0; i < len(appsName); i++ {
		appName := appsName[i]
		appDir := Conf.Appsdir + "/" + appName
		mkDir(appDir)
		// generage db password for every app
		Conf.MySQL.password = hashPass()

		cmd := exec.Command("docker", "cp", imgworkdir, appDir)
		logln("Execute: docker cp", imgworkdir, appDir)
		_, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}

		// Writing settings.php
		dstName := appDir + "/drupal/web/sites/default/settings.php"
		logln("Execute: cp -rf ", tplSettings, dstName)
		exec.Command("cp", "-rf", tplSettings, dstName).Run()

		logln("Write to file: ", dstName)
		textSlice := generateSettings(appName)
		writeFileln(dstName, textSlice)

		// Create files directory
		filesDir := appDir + "/drupal/web/sites/default/files"
		mkDir(filesDir)

		logln("Execute: cp -rf ", tplCompose, appDir)
		exec.Command("cp", "-rf", tplCompose, appDir+"/docker-compose.yml").Run()

		// logln("Execute: touch  ", tplCompose, appDir)
		// exec.Command("cp", "-rf", tplCompose, appDir + "/.env").Run()

		envs := generateDockEnv(appName)
		if len(envs) > 0 {
			writeFileln(appDir+"/.env", envs)
		}

	}

	logln("Execute: docker rm -f ", containerid)
	_, err = exec.Command("docker", "rm", "-f", containerid).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	logln("----- Done ------")
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
