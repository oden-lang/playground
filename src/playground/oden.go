package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

var odenc string

func findOdenc() {
	odenc = os.Getenv("ODENC")
	if odenc == "" {
		fmt.Println("Using odenc from $PATH")
		odenc = "odenc"
	}
}

func getOdenVersion() (string, error) {
	out, err := run(odenc, "--version")
	if err != nil {
		return "", errors.New("Failed to get odenc version: " + err.Error())
	}
	return strings.SplitAfter(out, " ")[0], nil
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

	out, err := run(odenc, "-p", tmpDir, "-o", tmpDir)
	if err != nil {
		return "", err
	}
	fmt.Println("Oden output:", out)

	goOutputPath := path.Join(tmpDir, "src", "main", "oden_out.go")
	if _, err = exec.LookPath("gofmt"); err == nil {
		if _, err = run("gofmt", "-w", goOutputPath); err != nil {
			return "", errors.New("Failed to gofmt Go output file: " + err.Error())
		}
	}
	goCode, err := ioutil.ReadFile(goOutputPath)
	if err != nil {
		return "", err
	}
	return string(goCode), nil
}