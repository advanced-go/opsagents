package logistics1

import (
	"context"
	"github.com/advanced-go/operations/landscape1"
	"github.com/advanced-go/stdlib/core"
	"net/url"
	"time"
)

const (
	queryDuration = time.Second * 2
)

type landscape struct {
	query func(region string) ([]landscape1.Entry, *core.Status)
}

func newLandscape() *landscape {
	return &landscape{
		query: func(region string) ([]landscape1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryDuration)
			defer cancel()
			values := make(url.Values)
			values.Add(landscape1.AssignedRegionKey, region)
			values.Add(landscape1.StatusKey, landscape1.StatusActive)
			return landscape1.Get(ctx, nil, values)
		},
	}
}
