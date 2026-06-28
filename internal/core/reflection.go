package core

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// BlueprintModel holds the metadata extracted from a Go struct, used by the core engine to generate dynamic handlers.
type BlueprintModel struct {
	Name      string
	Fields    []BlueprintField
	Prototype interface{}
}

// BlueprintField represents the metadata for a single field within a model, including its type and database constraints.
type BlueprintField struct {
	Name string
	Type string
	Kind reflect.Kind

	// Database constraints
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	Nullable      bool
	Size          int
	Default       string
	Searchable    bool
}

// NewInstance creates a new, zero-valued pointer to the model type defined in the blueprint.
func (bp *BlueprintModel) NewInstance() any {
	return reflect.New(
		reflect.TypeOf(bp.Prototype).Elem(),
	).Interface()
}

// BuildBlueprint analyzes a Go struct using reflection to extract field names, types, and custom 'crud' tags.
func BuildBlueprint(model interface{}) (*BlueprintModel, error) {
	t := reflect.TypeOf(model)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct")
	}

	bp := &BlueprintModel{
		Name:      t.Name(),
		Prototype: reflect.New(t).Interface(),
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		field := BlueprintField{
			Name: f.Name,
			Type: f.Type.Name(),
			Kind: f.Type.Kind(),
		}

		dbTag := f.Tag.Get("crud")
		if dbTag != "" {
			err := parseCrudTag(dbTag, &field)
			if err != nil {
				return nil, err
			}
		}

		bp.Fields = append(bp.Fields, field)
	}

	return bp, nil
}

// parseCrudTag interprets the comma-separated values within a 'crud' struct tag and updates the field metadata.
func parseCrudTag(tag string, field *BlueprintField) error {
	var err error
	parts := strings.SplitSeq(tag, ",")

	for part := range parts {
		switch {
		case part == "pk":
			field.PrimaryKey = true
		case part == "autoincrement":
			field.AutoIncrement = true
		case part == "unique":
			field.Unique = true
		case part == "nullable":
			field.Nullable = true
		case part == "searchable":
			field.Searchable = true
		case strings.HasPrefix(part, "size"):
			field.Size, err = strconv.Atoi(part[5:])
			if err != nil {
				return fmt.Errorf("invalid size: %w", err)
			}
		case strings.HasPrefix(part, "default"):
			field.Default = part[7:]
		default:
			return fmt.Errorf("unknown tag: %s", part)
		}
	}

	return nil
}

// GetSearchableFields returns a list of Go struct field names that are marked as searchable.
func (bp *BlueprintModel) GetSearchableFields() []string {
	searchableFields := make([]string, 0, len(bp.Fields))
	for _, field := range bp.Fields {
		if field.Searchable {
			searchableFields = append(searchableFields, field.Name)
		}
	}
	return searchableFields
}

// GetPrimaryKeyField returns the metadata for the field marked as the primary key.
func (bp *BlueprintModel) GetPrimaryKeyField() BlueprintField {
	var pkField BlueprintField
	for _, field := range bp.Fields {
		if field.PrimaryKey {
			pkField = field
			break
		}
	}
	return pkField
}
