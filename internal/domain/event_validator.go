package domain

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type EventValidator interface {
	Validate(event *Event) error
}

type JSONSchemaValidator struct {
	schemas  map[string]*gojsonschema.Schema
	registry EventRegistryInterface
}
type EventRegistryInterface interface {
	ResolveChannel(channel string) (string, string, error)
}

func NewJSONSchemaValidator(schemaDir string, registry EventRegistryInterface) (*JSONSchemaValidator, error) {
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
				Path:   "/" + filepath.ToSlash(absPath),
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
	return &JSONSchemaValidator{
		schemas:  schemas,
		registry: registry,
	}, nil
}

func (v *JSONSchemaValidator) Validate(event *Event) error {
	_, schemaName, err := v.registry.ResolveChannel(event.Type)
	schema, ok := v.schemas[schemaName]
	if !ok {
		return fmt.Errorf("no schema '%s' for event type: %s", schemaName, event.Type)
	}

	loader := gojsonschema.NewGoLoader(event.Payload)
	result, err := schema.Validate(loader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		var reasons []string
		for _, err := range result.Errors() {
			reasons = append(reasons, err.String())
		}
		return NewEventValidationError(strings.Join(reasons, "; "))
	}
	return nil
}

// func (v *JSONSchemaValidator) ReloadSchemas(schemaDir string) error {
// }
