package openapi

import (
	"github.com/morozovcookie/ogen/internal/location"
	"github.com/morozovcookie/ogen/jsonschema"
)

// Example is an OpenAPI Example.
type Example struct {
	Ref Ref

	Summary       string
	Description   string
	Value         jsonschema.Example
	ExternalValue string

	location.Pointer `json:"-" yaml:"-"`
}
