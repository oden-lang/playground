package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var odenc string

func init() {
	odenc = os.Getenv("ODENC")
	if odenc == "" {
		fmt.Println("Using odenc from $PATH")
		odenc = "odenc"
	}
	goRoot := os.Getenv("GOROOT")

	if goRoot == "" {
		fmt.Println("Setting GOROOT to /app/user/go for Heroku compatibility")
		os.Setenv("GOROOT", "/app/user/go")

		fmt.Println("Including /app/user/go/bin in PATH")
		path := os.Getenv("PATH")
		os.Setenv("PATH", path+":/app/user/go/bin")
	}
}

func getOdenVersion() (string, error) {
	out, err := run(odenc, "--version")
	if err != nil {
		return "", errors.New("Failed to get odenc version: " + err.Error())
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

	out, err := run(odenc, "-p"+tmpDir, "-o"+tmpDir, "-M")
	if err != nil {
		return "", err
	}
	fmt.Println("Oden output:", out)

	goOutputPath := path.Join(tmpDir, "src", "main.go")
	if _, err = run("gofmt", "-w", goOutputPath); err != nil {
		return "", errors.New("Failed to gofmt Go output file: " + err.Error())
	}
	goCode, err := ioutil.ReadFile(goOutputPath)
	if err != nil {
		return "", err
	}
	return string(goCode), nil
}
