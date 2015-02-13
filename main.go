package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type GMM struct{}

type moduleJson struct {
	Name         string
	Description  string
	Version      string
	Dependencies []string
	Author       string
}

type GmmInterface interface {
	getGoRoot() string
	getGoPath() string

	setGoPathTmp(path string) error

	execCmd(cmd string, wg *sync.WaitGroup) []byte
	install(name string)
	installDependencies()
	saveDependency(name string)

	successMessage(message string)
	warningMessage(message string)
	errorMessage(err error)
}

const (
	HEADER    = "\033[95m"
	OKBLUE    = "\033[94m"
	OKGREEN   = "\033[92m"
	WARNING   = "\033[93m"
	FAIL      = "\033[91m"
	ENDC      = "\033[0m"
	BOLD      = "\033[1m"
	UNDERLINE = "\033[4m"
)

func (g GMM) successMessage(message string) {
	fmt.Println(OKGREEN, message, ENDC)
}

func (g GMM) warningMessage(message string) {
	fmt.Println(WARNING, message, ENDC)
}

func (g GMM) errorMessage(err error) {
	fmt.Println(FAIL, err, ENDC)
}

func (g GMM) checkErr(err error) {
	if err != nil {
		g.errorMessage(err)
		os.Exit(1)
	}
}

func (g GMM) getGoPath() string {
	return os.Getenv("GOPATH")
}

func (g GMM) getGoRoot() string {
	return os.Getenv("GOROOT")
}

func (g GMM) setGoPathTmp(path string) error {
	err := os.Setenv("GOPATH", path)
	g.checkErr(err)
	return err
}

func (g GMM) execCmd(cmd string, wg *sync.WaitGroup) []byte {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	g.checkErr(err)

	wg.Done()
	return out
}

func (g GMM) install(name string) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	g.execCmd("go get "+name, wg)
	wg.Wait()
}

func (g GMM) installDependencies() {
	var module moduleJson

	file, err := ioutil.ReadFile("module.json")
	g.checkErr(err)

	err = json.Unmarshal(file, &module)
	g.checkErr(err)

	for _, el := range module.Dependencies {
		g.install(el)
	}
}

func (g GMM) saveDependency(name string) {
	var module moduleJson

	file, err := ioutil.ReadFile("module.json")
	g.checkErr(err)

	err = json.Unmarshal(file, &module)
	g.checkErr(err)

	module.Dependencies = append(module.Dependencies, name)

	moduleByte, err := json.MarshalIndent(module, "", "\t")
	g.checkErr(err)

	err = ioutil.WriteFile("module.json", moduleByte, 0644)
	g.checkErr(err)
}

func main() {

}
