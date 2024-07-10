package caseofficer1

import (
	"context"
	"github.com/advanced-go/operations/activity1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	duration = time.Second * 2
)

func activityLog(agentId string, content any) *core.Status {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	return activity1.Log(ctx, agentId, content)
}
