package logistics1

import (
	"context"
	"github.com/advanced-go/operations/activity1"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	logDuration = time.Second * 2
)

type operations struct {
	log func(handler messaging.OpsAgent, content any)
}

func newOperations() *operations {
	return &operations{
		log: func(handler messaging.OpsAgent, content any) {
			ctx, cancel := context.WithTimeout(context.Background(), logDuration)
			defer cancel()
			status := activity1.Log(ctx, handler.Uri(), content)
			if !status.OK() {
				handler.Handle(status, "")
			}
		},
	}
}
