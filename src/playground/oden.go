package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var oden string

func init() {
	oden = os.Getenv("ODEN_CLI")
	if oden == "" {
		fmt.Println("Using oden from $PATH")
		oden = "oden"
	}
	goRoot := os.Getenv("GOROOT")

	if goRoot == "" {
		fmt.Println("Setting GOROOT to /app/user/go for Heroku compatibility")
		os.Setenv("GOROOT", "/app/user/go")

		fmt.Println("Including /app/user/go/bin in PATH")
		path := os.Getenv("PATH")
		os.Setenv("PATH", path+":/app/user/go/bin")
	}

	os.Setenv("LC_ALL", "en_US.UTF-8")
	os.Setenv("LANG", "en_US.UTF-8")
	os.Setenv("LANGUAGE", "en_US.UTF-8")
}

func getOdenVersion() (string, error) {
	out, err := run(oden, "--version")
	if err != nil {
		return "", errors.New("Failed to get oden version: " + err.Error())
	}
	return strings.SplitAfter(strings.TrimSpace(out), " ")[0], nil
}

func compile(code string) (string, error) {
	tmpDir, err := ioutil.TempDir("", "oden")
	if err != nil {
		return "", errors.New("Failed to create temporary compile dir: " + err.Error())
	}

	err = os.MkdirAll(path.Join(tmpDir, "src"), 0775)
	if err != nil {
		return "", errors.New("Failed to create Oden src dir: " + err.Error())
	}
	defer os.Remove(tmpDir)

	odenSrc := path.Join(tmpDir, "src", "main.oden")
	err = ioutil.WriteFile(odenSrc, []byte(code), 0775)
	if err != nil {
		return "", errors.New("Failed to write Oden source file: " + err.Error())
	}

	out, err := run(oden, "-p"+tmpDir, "-o"+tmpDir, "-M", "build")
	if err != nil {
		return "", err
	}
	fmt.Println("Oden output:", out)

	goOutputPath := path.Join(tmpDir, "src", "main.go")
	if _, err := os.Stat(goOutputPath); os.IsNotExist(err) {
		return "", errors.New("The program must be the main package")
	}
	if _, err = run("gofmt", "-w", goOutputPath); err != nil {
		return "", errors.New("Failed to gofmt Go output file: " + err.Error())
	}
	goCode, err := ioutil.ReadFile(goOutputPath)
	if err != nil {
		return "", err
	}
	return string(goCode), nil
}
