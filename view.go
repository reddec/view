package view

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

// Layout base file name.
const layoutName = "_layout.gohtml"

// Load is the same as LoadTemplate but using new empty template as root.
func Load(store fs.FS, view string) (*template.Template, error) {
	return LoadTemplate(template.New(""), store, view)
}

// Load a single template (view) and all associated layouts (_layout.gohtml), starting from the top directory up to the current one.
//
// When using embedded FS (//go:embed), remember to use all suffixes (all:<pattern>) to ensure that files with underscores
// get embedded by the Go compiler. See https://github.com/golang/go/commit/36dbf7f7e63f3738795bb04593c3c011e987d1f3
func LoadTemplate(root *template.Template, store fs.FS, view string) (*template.Template, error) {
	dirs := strings.Split(strings.Trim(view, "/"), "/")
	dirs = dirs[:len(dirs)-1] // last segment is view itself
	// parse layouts from all dirs starting from the top and until the current dir
	for i := range dirs {
		fpath := path.Join(path.Join(dirs[:i+1]...), layoutName)
		content, err := fs.ReadFile(store, fpath)
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, fs.ErrNotExist) {
			continue // layout does not exists - skipping
		}
		if err != nil {
			return nil, fmt.Errorf("read layouad %q: %w", fpath, err)
		}

		child, err := root.Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", fpath, err)
		}
		root = child
	}
	// parse view it self
	content, err := fs.ReadFile(store, view)
	if err != nil {
		return nil, fmt.Errorf("parse view %q: %w", view, err)
	}
	return root.Parse(string(content))
}

// NewTemplate creates new [View] with provided root template and type-safe parameter.
func NewTemplate[T any](root *template.Template, store fs.FS, view string) (*View[T], error) {
	t, err := LoadTemplate(root, store, view)
	if err != nil {
		return nil, err
	}
	return &View[T]{
		parsed: t,
	}, nil
}

// New creates new [View] with empty root template and type-safe parameter.
func New[T any](store fs.FS, view string) (*View[T], error) {
	return NewTemplate[T](template.New(""), store, view)
}

// Must is convinient helper for wrapping around New* constructors.
// The function will panic if error parameter is not nil.
func Must[T any](v *View[T], e error) *View[T] {
	if e != nil {
		panic(e)
	}
	return v
}

// View is type-safe tiny wrapper around standard template.
type View[T any] struct {
	parsed *template.Template
}

// Render template as web page and content-type to text/html.
// It doesn't change response code. In case of error, some part of template
// could be sent to the client.
func (v *View[T]) Render(writer http.ResponseWriter, value T) error {
	writer.Header().Set("Content-Type", "text/html")
	return v.Execute(writer, value)
}

// Execute template and render content to the writer. It's just a type-safe wrapper around template.Execute.
func (v *View[T]) Execute(writer io.Writer, value T) error {
	return v.parsed.Execute(writer, value)
}

// Bytes result of template execution.
func (v *View[T]) Bytes(value T) ([]byte, error) {
	var buf bytes.Buffer
	err := v.Execute(&buf, value)
	return buf.Bytes(), err
}
