package core

import "context"

const (
	TriggerHook   = "@hook"
	TriggerManual = "@manual"
	TriggerCron   = "@cron"
)

type TriggerService interface {
	Trigger(context.Context, *Repository, *Hook) (*Build, error)
}
