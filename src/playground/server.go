package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

const defaultProgram = `(pkg main)

(import fmt)

(define (main) (fmt.Println "Hello, World!"))`

type RunRequest struct {
	OdenSource string `json:"odenSource"`
}
type RunResponse struct {
	Error         string        `json:"error"`
	GoOutput      string        `json:"goOutput"`
	ConsoleOutput *PlayResponse `json:"consoleOutput"`
}

type ViewModel struct {
	Version    string
	OdenSource string
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
		})
	}).Methods("GET")

	router.HandleFunc("/program/{prg}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		prg := vars["prg"]
		data, err := base64.StdEncoding.DecodeString(prg)
		if err != nil {
			fmt.Println(string(data))
			fmt.Fprintf(os.Stderr, "Failed to decode program: %s\n", err)
			r.HTML(w, http.StatusBadRequest, "invalid-program", nil)
			return
		}
		r.HTML(w, http.StatusOK, "index", ViewModel{
			version,
			string(data),
		})
	}).Methods("GET")

	router.HandleFunc("/compile", func(w http.ResponseWriter, req *http.Request) {

		var runReq RunRequest
		if err := json.NewDecoder(req.Body).Decode(&runReq); err != nil {
			http.Error(w, "Could not decode JSON", http.StatusBadRequest)
			return
		}

		goCode, err := compile(runReq.OdenSource)
		if err != nil {
			fmt.Println("Failed to compile due to", err, ":\n", runReq.OdenSource)
			r.JSON(w, http.StatusOK, RunResponse{
				err.Error(),
				"",
				nil,
			})
			return
		}

		consoleOutput, err := runGoPkg(goCode, version)
		if err != nil {
			fmt.Println("Failed to run due to", err, ":\n", goCode)
			r.JSON(w, http.StatusOK, RunResponse{
				err.Error(),
				"",
				nil,
			})
			return
		}
		if consoleOutput.Errors != "" {
			fmt.Println("Run with errors:\n", consoleOutput.Errors)
		} else {
			fmt.Println("Run:\n", goCode)
		}

		r.JSON(w, http.StatusOK, RunResponse{
			"",
			goCode,
			consoleOutput,
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
