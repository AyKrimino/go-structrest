package core

import "context"

// HookType represents the different lifecycle events available for models.
type HookType string

const (
	BeforeCreate HookType = "before_create"
	AfterCreate  HookType = "after_create"
	BeforeUpdate HookType = "before_update"
	AfterUpdate  HookType = "after_update"
	BeforeDelete HookType = "before_delete"
	AfterDelete  HookType = "after_delete"
)

// BeforeCreateHook is an optional interface that models can implement to execute logic
// before a new record is created in the database.
type BeforeCreateHook interface {
	BeforeCreate(ctx context.Context) error
}

// AfterCreateHook is an optional interface that models can implement to execute logic
// after a new record has been successfully created.
type AfterCreateHook interface {
	AfterCreate(ctx context.Context) error
}

// BeforeUpdateHook is an optional interface that models can implement to execute logic
// before an existing record is updated.
type BeforeUpdateHook interface {
	BeforeUpdate(ctx context.Context) error
}

// AfterUpdateHook is an optional interface that models can implement to execute logic
// after an existing record has been successfully updated.
type AfterUpdateHook interface {
	AfterUpdate(ctx context.Context) error
}

// BeforeDeleteHook is an optional interface that models can implement to execute logic
// before a record is deleted. Returning an error here will abort the deletion.
type BeforeDeleteHook interface {
	BeforeDelete(ctx context.Context) error
}

// AfterDeleteHook is an optional interface that models can implement to execute logic
// after a record has been successfully deleted.
type AfterDeleteHook interface {
	AfterDelete(ctx context.Context) error
}
