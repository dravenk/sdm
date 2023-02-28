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
	tplCompose  = "docker-compose.yaml"
	appPath     = "drupal"
)

var cmdInput string

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(removeCmd)

	loadConf()
}

func main() {
	rootCmd.Execute()

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
			logln(err)
			return
		}
		containerid := strings.TrimSuffix(string(cid), "\n")
		Conf.ContainerId = containerid

		defer removeContainer(containerid)

		imgworkdir := containerid + `:` + Conf.Workdir
		for i := 0; i < len(appsName); i++ {
			appName := appsName[i]
			appDir := filepath.Join(Conf.Appsdir, appName)
			initProjectFiles(appName, appDir, imgworkdir)
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

func removeContainer(containerid string) bool {
	logln("Execute: docker rm -f ", containerid)
	if _, err := exec.Command("docker", "rm", "-f", containerid).CombinedOutput(); err != nil {
		logln(err)
		return false
	}
	return true
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
	dstName := pathJoin(appDir, appPath, "web/sites/default/settings.php")
	logln("Execute: cp -rf ", tplSettings, dstName)
	exec.Command("cp", "-rf", tplSettings, dstName).Run()

	logln("Write to file: ", dstName)
	writeFileln(dstName, generateSettings(appName))

	// Create files directory
	filesDir := pathJoin(appDir, appPath, "web/sites/default/files")
	mkDir(filesDir, os.ModePerm)
	os.Chmod(filesDir, os.ModePerm)

	logln("Execute: cp -rf ", tplCompose, appDir)
	exec.Command("cp", "-rf", tplCompose, pathJoin(appDir, tplCompose)).Run()
	writeFileln(pathJoin(appDir, ".env"), generateDockEnv(appName))
}

func startUpApps(appName, appDir string) {
	logln("Execute: docker-compose", "-f", pathJoin(appDir, tplCompose), InputUp, "-d")
	exec.Command("docker-compose", "-f", pathJoin(appDir, tplCompose), InputUp, "-d").Run()
}

func downApps(appName, appDir string) {
	logln("Execute: docker-compose", "-f", pathJoin(appDir, tplCompose), InputDown)
	exec.Command("docker-compose", "-f", pathJoin(appDir, tplCompose), InputDown).Run()
}

func removeApps() {
	logln("Execute: RemoveAll ", Conf.Appsdir)
	if err := os.RemoveAll(Conf.Appsdir); err != nil {
		log.Fatal(err)
	}
}
