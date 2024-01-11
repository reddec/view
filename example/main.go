package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/reddec/view"
)

//go:embed all:views
var views embed.FS

func main() {
	index := view.Must(view.New[string](views, "views/index.gohtml"))
	about := view.Must(view.New[string](views, "views/info/about.gohtml"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index.Render(w, "the index page")
	})
	http.HandleFunc("/info/about", func(w http.ResponseWriter, r *http.Request) {
		about.Render(w, "made by RedDec")
	})
	fmt.Println("ready on :8080")
	panic(http.ListenAndServe(":8080", nil))
}
