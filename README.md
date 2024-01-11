# View [![GoDoc](https://godoc.org/github.com/reddec/view?status.png)](https://godoc.org/github.com/reddec/view)

The `view` package provides a type-safe and hierarchical (layouts) way to load and render Go HTML templates using the standard library's [`html/template`](https://pkg.go.dev/html/template) package. It supports loading templates from any sources exposed as [`fs.FS`](https://pkg.go.dev/io/fs#FS). The package comes with no external dependencies and is designed for use in web applications.

This is extremly light library based on gist https://gist.github.com/reddec/312367d75cc03f1ee49bae74c52a6b31 and has zero external dependecies.

Key points:

- **Hierarchical**: The templates are loaded in a hierarchical way, allowing you to have a base layout and extend it with partials or views at different levels. Layouts  defined in each directory as `_layout.gohtml` file and can be extended.
- **Type-safe**: The package provides a type-safe wrapper around the standard `html/template` library using a custom `View` struct.


## [Example](example/)

Layout

```
├── main.go
└── views
    ├── _layout.gohtml
    ├── index.gohtml
    └── info
        ├── _layout.gohtml
        └── about.gohtml
```


And the code

```go
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
```

- note: `all:view` - the `all:` prefix is required in order to include files with underscore in name prefix

## Installation

```bash
go get github.com/reddec/view
```