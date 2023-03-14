package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func init() {
	loadConf()
}

func main() {
	// cli
	Execute()

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

		wg := new(sync.WaitGroup)
		wg.Add(len(appsName))

		imgworkdir := containerid + `:` + Conf.Workdir
		for i := 0; i < len(appsName); i++ {
			appName := appsName[i]
			appDir := filepath.Join(Conf.Appsdir, appName)
			go initProjectFiles(wg, appName, appDir, imgworkdir)
		}
		wg.Wait()
	}
	// if cmdInput == InputRemove {
	// 	removeApps()
	// }

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

func initProjectFiles(wg *sync.WaitGroup, appName, appDir, imgworkdir string) {
	defer wg.Done()

	mkDir(appDir, os.ModePerm)
	Conf.MySQL.password = hashPass()
	appDirPath := filepath.Join(appDir,filepath.Base(Conf.Workdir))
	
	if !isNotExist(appDirPath) {
		return
	}

	logln("Execute: docker cp", imgworkdir, appDir)
	cmdcp := exec.Command("docker", "cp", imgworkdir, appDir)
	if _, err := cmdcp.CombinedOutput(); err != nil {
		log.Fatal(err)
	}

	// Writing settings.php
	dstName := pathJoin(appDir, appPath, "web/sites/default/settings.php")
	srcName := pathJoin(projectPath, tplSettings)
	logln("Execute: cp -rf ", projectPath, dstName)
	exec.Command("cp", "-rf", srcName , dstName).Run()

	logln("Write to file: ", dstName)
	writeFileln(dstName, generateSettings(appName))

	// Create files directory
	filesDir := pathJoin(appDir, appPath, "web/sites/default/files")
	mkDir(filesDir, os.ModePerm)
	os.Chmod(filesDir, os.ModePerm)

	srcName = pathJoin(projectPath, tplCompose)
	logln("Execute: cp -rf ", srcName, appDir)
	exec.Command("cp", "-rf", srcName, pathJoin(appDir, tplCompose)).Run()
	writeFileln(pathJoin(appDir, ".env"), generateDockEnv(appName))
}

func startUpApps(appName, appDir string) {
	srcName :=  pathJoin(appDir, tplCompose)
	logln("Execute: docker-compose", "-f",srcName, InputUp, "-d")
	exec.Command("docker-compose", "-f", srcName, InputUp, "-d").Run()
}

func downApps(appName, appDir string) {
	srcName :=  pathJoin(appDir, tplCompose)
	logln("Execute: docker-compose", "-f",srcName, InputDown)
	exec.Command("docker-compose", "-f", srcName, InputDown).Run()
}

func removeApps() {
	logln("Execute: RemoveAll ", Conf.Appsdir)
	if err := os.RemoveAll(Conf.Appsdir); err != nil {
		log.Fatal(err)
	}
}
