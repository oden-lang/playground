package main

import (
	"net/http"
	"os"

	"github.com/oden-lang/playground/Godeps/_workspace/src/github.com/codegangsta/negroni"
	"github.com/oden-lang/playground/Godeps/_workspace/src/github.com/unrolled/render"
)

type ViewModel struct {
	OdenSource    string
	GoOutput      string
	ConsoleOutput string
}

func main() {
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
