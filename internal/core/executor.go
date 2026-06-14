package core

import "context"

func RunHook(ctx context.Context, model any, hookType HookType) error {
	switch hookType {
	case BeforeCreate:
		if hook, ok := model.(BeforeCreateHook); ok {
			return hook.BeforeCreate(ctx)
		}
	case AfterCreate:
		if hook, ok := model.(AfterCreateHook); ok {
			return hook.AfterCreate(ctx)
		}
	case BeforeUpdate:
		if hook, ok := model.(BeforeUpdateHook); ok {
			return hook.BeforeUpdate(ctx)
		}
	case AfterUpdate:
		if hook, ok := model.(AfterUpdateHook); ok {
			return hook.AfterUpdate(ctx)
		}
	case BeforeDelete:
		if hook, ok := model.(BeforeDeleteHook); ok {
			return hook.BeforeDelete(ctx)
		}
	case AfterDelete:
		if hook, ok := model.(AfterDeleteHook); ok {
			return hook.AfterDelete(ctx)
		}
	}
	return nil
}
