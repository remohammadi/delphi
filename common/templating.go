package common

import (
	"fmt"
	"html/template"
	"io"
	"path"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/oxtoacart/bpool"
)

var (
	templates map[string]*template.Template
	bufpool   *bpool.BufferPool
)

// Load templates on program initialisation
func init() {
	bufpool = bpool.NewBufferPool(64)

	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templatesDir := ConfigString("TEMPLATES_DIR")

	layoutsPath := path.Join(templatesDir, "layouts", "*.tmpl")
	layouts, err := filepath.Glob(layoutsPath)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed while loading templates inside %s\n", layoutsPath)
	}

	includesPath := path.Join(templatesDir, "includes", "*.tmpl")
	includes, err := filepath.Glob(includesPath)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed while loading templates inside %s\n", includesPath)
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, layout := range layouts {
		files := append(includes, layout)
		name := filepath.Base(layout)
		templates[name] = template.Must(template.ParseFiles(files...))
	}
}

// RenderTemplate is a wrapper around template.ExecuteTemplate.
func RenderTemplate(w io.Writer, base, name string, data map[string]interface{}) (func(), error) {
	// Ensure the template exists in the map.
	tmpl, ok := templates[name]
	if !ok {
		return nil, fmt.Errorf("The template %s does not exist.", name)
	}

	// Create a buffer to temporarily write to and check if any errors were encounted.
	buf := bufpool.Get()

	context := map[string]interface{}{
		"DOMAIN": ConfigString("SERVER_NAME"),
	}
	for k, v := range data {
		context[k] = v
	}

	err := tmpl.ExecuteTemplate(buf, base, context)
	if err != nil {
		bufpool.Put(buf)
		return nil, err
	}

	return func() {
		defer bufpool.Put(buf)
		buf.WriteTo(w)
	}, nil
}
