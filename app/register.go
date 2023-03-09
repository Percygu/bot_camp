package app

import "context"

func RegisterTask(ctx context.Context) {
	TriggerWorkshop(ctx)
	RemindTalk(ctx)
}
