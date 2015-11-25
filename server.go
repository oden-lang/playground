package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/oden-lang/playground/Godeps/_workspace/src/github.com/codegangsta/negroni"
	"github.com/oden-lang/playground/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/oden-lang/playground/Godeps/_workspace/src/github.com/unrolled/render"
)

const defaultProgram = `(pkg main)

(import fmt)

(define (main) (fmt.Println "Hello, World!"))`

type ViewModel struct {
	Version          string
	OdenSource       string
	GoOutput         string
	ConsoleOutput    string
	CompilationError error
}

func main() {
	findOdenc()
	version, err := getOdenVersion()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
		return
	}
	fmt.Println("Oden version:", version)

	r := render.New(render.Options{
		Layout:     "layout",
		Extensions: []string{".html"},
	})
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "index", ViewModel{
			version,
			defaultProgram,
			"",
			"",
			nil,
		})
	}).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			r.HTML(w, http.StatusBadRequest, "bad-request", nil)
			return
		}

		source := req.FormValue("odenSource")
		goCode, err := compile(source)

		if err != nil {
			fmt.Println("Failed to compile:", err)
			r.HTML(w, http.StatusOK, "index", ViewModel{
				version,
				source,
				"",
				"",
				err,
			})
			return
		}

		r.HTML(w, http.StatusOK, "index", ViewModel{
			version,
			source,
			goCode,
			"",
			nil,
		})
	}).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	n := negroni.Classic()
	n.UseHandler(router)
	n.Use(negroni.NewStatic(http.Dir("public")))
	n.Run(":" + port)
}
