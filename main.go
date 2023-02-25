package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	tplSettings = "settings.php"
	tplCompose  = "docker-compose.yml"
)

var cmdInput string

func init() {
	cmdInput = initCommand()
	loadConf()
}

func main() {

	appsName := Conf.AppsName
	if len(appsName) == 0 {
		logln("Exit. Not found any applications name in configuration.", Conf.Image)
		return
	}

	if cmdInput == InputInit {
		mkDir(Conf.Appsdir, os.ModePerm)

		// Create container: docker create drupal:latest
		logln("Execute: docker create", Conf.Image)
		cmdcid := exec.Command("docker", "create", Conf.Image)
		cid, err := cmdcid.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		containerid := strings.TrimSuffix(string(cid), "\n")
		imgworkdir := containerid + `:` + Conf.Workdir
		for i := 0; i < len(appsName); i++ {
			appName := appsName[i]
			appDir := filepath.Join(Conf.Appsdir, appName)
			initProjectFiles(appName, appDir, imgworkdir)
		}
		logln("Execute: docker rm -f ", containerid)
		if _, err = exec.Command("docker", "rm", "-f", containerid).CombinedOutput(); err != nil {
			log.Fatal(err)
		}
	}
	if cmdInput == InputRemove {
		removeApps()
	}

	for i := 0; i < len(appsName); i++ {
		appName := appsName[i]
		appDir := filepath.Join(Conf.Appsdir, appName)
		switch cmdInput {
		case InputUp:
			startUpApps(appName, appDir)
		case InputDown:
			downApps(appName, appDir)
		}
	}
}

func initProjectFiles(appName, appDir, imgworkdir string) {
	mkDir(appDir, os.ModePerm)
	// generage db password for every app
	Conf.MySQL.password = hashPass()

	logln("Execute: docker cp", imgworkdir, appDir)
	cmdcp := exec.Command("docker", "cp", imgworkdir, appDir)
	if _, err := cmdcp.CombinedOutput(); err != nil {
		log.Fatal(err)
	}

	// Writing settings.php
	dstName := appDir + "/drupal/web/sites/default/settings.php"
	logln("Execute: cp -rf ", tplSettings, dstName)
	exec.Command("cp", "-rf", tplSettings, dstName).Run()

	logln("Write to file: ", dstName)
	writeFileln(dstName, generateSettings(appName))

	// Create files directory
	filesDir := appDir + "/drupal/web/sites/default/files"
	mkDir(filesDir, os.ModePerm)
	os.Chmod(filesDir, os.ModePerm)

	logln("Execute: cp -rf ", tplCompose, appDir)
	exec.Command("cp", "-rf", tplCompose, appDir+"/docker-compose.yml").Run()
	writeFileln(appDir+"/.env", generateDockEnv(appName))
}

func startUpApps(appName, appDir string) {
	logln("Execute: docker-compose", "-f", appDir+"/docker-compose.yml", InputUp, "-d")
	exec.Command("docker-compose", "-f", appDir+"/docker-compose.yml", InputUp, "-d").Run()
}

func downApps(appName, appDir string) {
	logln("Execute: docker-compose", "-f", appDir+"/docker-compose.yml", InputDown)
	exec.Command("docker-compose", "-f", appDir+"/docker-compose.yml", InputDown).Run()
}

func removeApps() {
	logln("Execute: RemoveAll ", Conf.Appsdir)
	// Remove all the directories and files
	if err := os.RemoveAll(Conf.Appsdir); err != nil {
		log.Fatal(err)
	}
}
