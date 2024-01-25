package autopilot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type FixtureTemplater struct {
	coreTemplate *template.Template
}

func (t *FixtureTemplater) ParseToBytes(contents string) (*bytes.Buffer, error) {
	clone, err := t.coreTemplate.Clone()
	if err != nil {
		return nil, fmt.Errorf("error cloning core template: %s", err)
	}
	target, err := clone.Parse(contents)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %s", err)
	}
	data := bytes.NewBuffer([]byte{})
	if err = target.Execute(data, nil); err != nil {
		return nil, err
	}
	return data, nil
}

func (t *FixtureTemplater) Use(contents string) (string, error) {
	data, err := t.ParseToBytes(contents)
	if err != nil {
		return "", err
	}
	output := bytes.NewBuffer([]byte{})
	err = json.Indent(output, data.Bytes(), "", "  ")
	if err != nil {
		return "", fmt.Errorf("error during json indenting: %s", err)
	}
	return output.String(), nil
}

func (t *FixtureTemplater) UseKey(key string) (string, error) {
	return t.Use(fmt.Sprintf("{{ template \"%s\" }}", key))
}

func (t *FixtureTemplater) UseFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading file '%s': %s", path, err)
	}
	return t.Use(string(data))
}

// NewFixtureTemplater creates a FixtureTemplater using the templates load from the passed in templateDir
//
// If no templateDir is passed in then ./testdata/templates is used
func NewFixtureTemplater(templateDirs ...string) *FixtureTemplater {
	if len(templateDirs) == 0 {
		templateDirs = append(templateDirs, "./testdata/templates")
	}
	output := FixtureTemplater{}
	// A map is used to act like a set for filepath dedupilcation
	files := map[string]bool{}
	for _, dir := range templateDirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files[path] = true
			}
			return nil
		})
		if err != nil {
			panic(fmt.Errorf("error during loading template files: %s", err))
		}
	}
	var templateFiles []string
	for k := range files {
		templateFiles = append(templateFiles, k)
	}
	tmpl, err := template.New("").Funcs(sprig.TxtFuncMap()).ParseFiles(templateFiles...)
	if err != nil {
		panic(fmt.Errorf("error during template initialization: %s", err))
	}
	output.coreTemplate = tmpl
	return &output
}
