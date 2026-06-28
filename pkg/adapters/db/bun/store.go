package bun_store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/AyKrimino/go-structrest/pkg/adapters/db"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/schema"
)

type BunStore struct {
	db *bun.DB
}

func NewBunStore(db *bun.DB) *BunStore {
	return &BunStore{
		db: db,
	}
}

func (s *BunStore) Create(ctx context.Context, model any) error {
	_, err := s.db.NewInsert().
		Model(model).
		Exec(ctx)
	return err
}

func (s *BunStore) FindByID(ctx context.Context, model any, id any) error {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		bunTag := field.Tag.Get("bun")

		isPk := false
		for part := range strings.SplitSeq(bunTag, ",") {
			if strings.TrimSpace(part) == "pk" {
				isPk = true
				break
			}
		}

		if isPk {
			fieldValue := v.Field(i)
			if fieldValue.CanSet() {
				idVal := reflect.ValueOf(id)
				if idVal.Type().ConvertibleTo(fieldValue.Type()) {
					fieldValue.Set(idVal.Convert(fieldValue.Type()))
				}
			}
			break
		}
	}

	err := s.db.NewSelect().
		Model(model).
		WherePK().
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return db.ErrResourceNotFound
	}
	return err
}

func (s *BunStore) FindAll(ctx context.Context, model any, opts db.QueryOptions) error {
	query := s.db.NewSelect().
		Model(model)

	if opts.Search != "" && len(opts.SearchableFields) > 0 {
		searchOp := "LIKE"
		if s.db.Dialect().Name() == dialect.PG {
			searchOp = "ILIKE"
		}

		query = query.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			for _, field := range opts.SearchableFields {
				sq = sq.WhereOr(fmt.Sprintf("? %s ?", searchOp), bun.Ident(field), "%"+opts.Search+"%")
			}
			return sq
		})
	}

	if opts.SortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", opts.SortBy, opts.Order))
	}

	query.Limit(opts.Limit).Offset(opts.Offset)

	return query.Scan(ctx)
}

func (s *BunStore) Update(ctx context.Context, model any) error {
	_, err := s.db.NewUpdate().
		Model(model).
		WherePK().
		Exec(ctx)
	return err
}

func (s *BunStore) Delete(ctx context.Context, model any) error {
	_, err := s.db.NewDelete().
		Model(model).
		WherePK().
		Exec(ctx)
	return err
}

func (s *BunStore) GetColumnName(model any, goFieldName string) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	tables := schema.NewTables(s.db.Dialect())
	tables.Register(reflect.New(t).Interface())

	table := tables.ByModel(t.Name())
	for _, f := range table.Fields {
		if f.GoName == goFieldName {
			return f.Name
		}
	}

	return goFieldName
}
