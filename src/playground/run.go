package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func run(prg string, params ...string) (string, error) {
	cmd := exec.Command(prg, params...)
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	var errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	err := cmd.Run()

	fmt.Println("out", out.String(), "err", errout.String())
	if err != nil {
		return "", errors.New(errout.String())
	}
	return out.String(), nil
}
