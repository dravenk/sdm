package main

import (
	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func writeFileln(dstFileName string, textSlice []string) {
	f, err := os.OpenFile(dstFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
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

func mkDir(dir string) {
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
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
