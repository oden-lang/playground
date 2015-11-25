package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/oden-lang/playground/Godeps/_workspace/src/github.com/codegangsta/negroni"
	"github.com/oden-lang/playground/Godeps/_workspace/src/github.com/unrolled/render"
)

func getOdenVersion() (string, error) {
	fmt.Println("PATH", os.Getenv("PATH"))
	cmd := exec.Command("odenc", "version")
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

type ViewModel struct {
	OdenSource    string
	GoOutput      string
	ConsoleOutput string
}

func main() {
	version, err := getOdenVersion()
	fmt.Println("Oden version:", version, err)

	r := render.New(render.Options{
		Layout:     "layout",
		Extensions: []string{".html"},
	})
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "index", ViewModel{
			"(pkg main)",
			"package main",
			"",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Use(negroni.NewStatic(http.Dir("public")))
	n.Run(":" + port)
}
