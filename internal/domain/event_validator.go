package domain

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

type EventValidator interface {
	Validate(event *Event) error
}

type JSONSchemaValidator struct {
	schemas map[string]*gojsonschema.Schema
}

func NewJSONSchemaValidator(schemaDir string) (*JSONSchemaValidator, error) {
	schemas := make(map[string]*gojsonschema.Schema)
	files, err := os.ReadDir(schemaDir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			schemaPath := filepath.Join(schemaDir, file.Name())
			absPath, err := filepath.Abs(schemaPath)
			if err != nil {
				return nil, err
			}
			u := &url.URL{
				Scheme: "file",
				Path:   "/" + filepath.ToSlash(absPath), // "/" for Windows
			}
			schemaLoader := gojsonschema.NewReferenceLoader(u.String())
			schema, err := gojsonschema.NewSchema(schemaLoader)
			if err != nil {
				return nil, fmt.Errorf("failed to load schema %s: %w", file.Name(), err)
			}
			schemaType := file.Name()[:len(file.Name())-len(".schema.json")]
			schemas[schemaType] = schema
		}
	}
	return &JSONSchemaValidator{schemas: schemas}, nil
}

func (v *JSONSchemaValidator) Validate(event *Event) error {
	schema, ok := v.schemas[event.Type]
	if !ok {
		return fmt.Errorf("no schema for event type: %s", event.Type)
	}
	loader := gojsonschema.NewGoLoader(event.Payload)
	result, err := schema.Validate(loader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return fmt.Errorf("validation error: %v", result.Errors())
	}
	return nil
}
