package core

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type BlueprintModel struct {
	Name      string
	Fields    []BlueprintField
	Prototype interface{}
}

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

type ValidationRule struct {
	Name  string
	Value string
}

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

func (bp *BlueprintModel) NewInstance() any {
	return reflect.New(
		reflect.TypeOf(bp.Prototype).Elem(),
	).Interface()
}

func (bp *BlueprintModel) GetSearchableFields() []string {
	searchableFields := make([]string, 0, len(bp.Fields))
	for _, field := range bp.Fields {
		if field.Searchable {
			searchableFields = append(searchableFields, field.Name)
		}
	}
	return searchableFields
}

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
