package core

import "context"

type HookType string

const (
	BeforeCreate HookType = "before_create"
	AfterCreate  HookType = "after_create"
	BeforeUpdate HookType = "before_update"
	AfterUpdate  HookType = "after_update"
	BeforeDelete HookType = "before_delete"
	AfterDelete  HookType = "after_delete"
)

type BeforeCreateHook interface {
	BeforeCreate(ctx context.Context) error
}

type AfterCreateHook interface {
	AfterCreate(ctx context.Context) error
}

type BeforeUpdateHook interface {
	BeforeUpdate(ctx context.Context) error
}

type AfterUpdateHook interface {
	AfterUpdate(ctx context.Context) error
}

type BeforeDeleteHook interface {
	BeforeDelete(ctx context.Context) error
}

type AfterDeleteHook interface {
	AfterDelete(ctx context.Context) error
}
