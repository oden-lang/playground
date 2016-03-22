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

const defaultProgram = `package main

import fmt

square : int -> int
square(n) = n * n

twice : forall a. (a -> a) -> a -> a
twice(f, x) = f(f(x))

main : -> ()
main() = fmt.Println("2\x2074 =", twice(square, 2))`

type CodeRequest struct {
	OdenSource string `json:"odenSource"`
}
type RunResponse struct {
	Error         string        `json:"error"`
	GoOutput      string        `json:"goOutput"`
	ConsoleOutput *PlayResponse `json:"consoleOutput"`
}
type SaveResponse struct {
	ProgramId string `json:"programId"`
	Path      string `json:"path"`
}

type ViewModel struct {
	Version    string
	OdenSource string
	Deprecated bool
}

func main() {
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
			false,
		})
	}).Methods("GET")

	router.HandleFunc("/p/{id}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]

		prg, err := findProgram(id)
		fmt.Println(prg, err)
		if err != nil {
			r.HTML(w, http.StatusInternalServerError, "500", nil)
			return
		} else if prg == nil {
			r.HTML(w, http.StatusNotFound, "404", nil)
			return
		}

		r.HTML(w, http.StatusOK, "index", ViewModel{
			version,
			*prg,
			false,
		})
	}).Methods("GET")

	router.HandleFunc("/compile", func(w http.ResponseWriter, req *http.Request) {

		var runReq CodeRequest
		if err := json.NewDecoder(req.Body).Decode(&runReq); err != nil {
			http.Error(w, "Could not decode JSON", http.StatusBadRequest)
			return
		}

		goCode, err := compile(runReq.OdenSource)
		if err != nil {
			fmt.Println("Failed to compile due to", err, ":\n", runReq.OdenSource)
			go mailer.SendOdenCompilationError(runReq.OdenSource, err.Error())
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
			go mailer.SendGoRunError(runReq.OdenSource, goCode, err.Error())
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
			go mailer.SendCodeSuccess(runReq.OdenSource, goCode, consoleOutput)
			fmt.Println("Run:\n", goCode)
		}

		r.JSON(w, http.StatusOK, RunResponse{
			"",
			goCode,
			consoleOutput,
		})
	}).Methods("POST")

	router.HandleFunc("/p", func(w http.ResponseWriter, req *http.Request) {
		var saveReq CodeRequest
		if err := json.NewDecoder(req.Body).Decode(&saveReq); err != nil {
			http.Error(w, "Could not decode JSON", http.StatusBadRequest)
			return
		}

		id, err := saveProgram(saveReq.OdenSource)
		if err != nil {
			fmt.Println("Failed to save code:", err)
			http.Error(w, "Could not save program", http.StatusInternalServerError)
			return
		}

		r.JSON(w, http.StatusOK, SaveResponse{
			id,
			"/p/" + id,
		})
	}).Methods("POST")

	// For backwards compatibility
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
			true,
		})
	}).Methods("GET")

	router.HandleFunc("/about", func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "about", ViewModel{
			version,
			"",
			true,
		})
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	n := negroni.Classic()
	n.UseHandler(router)
	n.Use(negroni.NewStatic(http.Dir("public")))
	n.Run(":" + port)
}
