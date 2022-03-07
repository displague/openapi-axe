package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/getkin/kin-openapi/openapi3"
	"sigs.k8s.io/yaml"
)

var (
	infile string
	dir    string

	schemaDir        = "components/schemas"
	examplesDir      = "components/examples"
	headersDir       = "components/headers"
	requestBodiesDir = "components/requestBodies"
	responsesDir     = "components/responses"
	parametersDir    = "components/parameters"
	pathsDir         = "paths"
	tagsDir          = "tags"

	filemode = fs.FileMode(0644)
)

func main() {
	flag.StringVar(&infile, "i", "", "The OpenAPI source to split")
	flag.StringVar(&dir, "d", "", "The target directory")
	flag.Parse()

	if infile == "" || dir == "" {
		flag.PrintDefaults()
		return
	}

	doc, err := openapi3.NewLoader().LoadFromFile(infile)
	if err != nil {
		panic(err)
	}

	for name, schema := range doc.Components.Schemas {
		if schema.Ref != "" {
			continue
		}
		filename := path.Join(".", schemaDir, name+".yaml")
		if err = writeObject(schema, path.Join(dir, filename)); err != nil {
			panic(err)
		}
		schema.Ref = filename
	}

	for name, p := range doc.Paths {
		if p.Ref != "" {
			continue
		}
		for method, op := range p.Operations {
			name = op.OperationID
			filename := path.Join(".", pathsDir, name+".yaml")
			if err = appendObject(p, path.Join(dir, filename)); err != nil {
				panic(err)
			}
			p.Ref = fmt.Sprintf(filename, "#/paths/%s", name)
		}
	}
}

func writeObject(object interface{}, filename string) error {
	data, err := yaml.Marshal(object)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, filemode)
}

func appendObject(object interface{}, filename string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, filemode)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	data, err := yaml.Marshal(object)
	if err != nil {
		return err
	}
	if _, err = f.Write([]byte(fmt.Sprintf("---\n%s\n", data))); err != nil {
		panic(err)
	}
}
