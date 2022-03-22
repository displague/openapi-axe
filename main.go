package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/jsonpointer"
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

	filemode = fs.FileMode(0o644)
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

	if err = os.MkdirAll(path.Join(dir, schemaDir), os.FileMode(0o1755)); err != nil {
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

	if err = os.MkdirAll(path.Join(dir, pathsDir), os.FileMode(0o1755)); err != nil {
		panic(err)
	}

	for name, p := range doc.Paths {
		if p.Ref != "" {
			continue
		}

		filename := path.Join(dir, pathsDir, jsonpointer.Escape(name)+".yaml")
		refName := path.Join(pathsDir, jsonpointer.Escape(name)+".yaml")
		if err = writeObject(p, filename); err != nil {
			panic(err)
		}
		refOp := &openapi3.PathItem{Ref: fmt.Sprintf("./%s", refName)}
		doc.Paths[name] = refOp
		/**
		for verb, op := range p.Operations() {
			opDir := path.Join(dir, pathsDir, jsonpointer.Escape(name))
			filename := strings.ToLower(verb) + ".yaml"
			refName := path.Join(pathsDir, jsonpointer.Escape(name), filename)
			if err = os.MkdirAll(opDir, os.FileMode(0o1755)); err != nil {
				panic(err)
			}
			if err = writeObject(op, path.Join(opDir, filename)); err != nil {
				panic(err)
			}
			refOp := &openapi3.Operation{Ref: fmt.Sprintf("#/%s", refName)}
			op = refOp
		}
		**/
	}

	if err = writeObject(doc, path.Join(dir, "openapi3.yaml")); err != nil {
		panic(err)
	}

	return
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
		return err
	}

	defer f.Close()
	data, err := yaml.Marshal(object)
	if err != nil {
		return err
	}
	if _, err = f.Write([]byte(fmt.Sprintf("---\n%s\n", data))); err != nil {
		return err
	}
	return nil
}
