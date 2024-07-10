package caseofficer1

import (
	"context"
	"errors"
	"github.com/advanced-go/operations/assignment1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	updateDuration    = time.Second * 2
	addStatusDuration = time.Second * 2
)

type assignment struct {
	update    func(agentId string, origin core.Origin) ([]assignment1.Entry, *core.Status)
	addStatus func(msg *messaging.Message) *core.Status
}

func newAssignment() *assignment {
	return &assignment{
		update: func(agentId string, origin core.Origin) ([]assignment1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), updateDuration)
			defer cancel()
			return assignment1.Update(ctx, agentId, origin)
		},
		addStatus: func(msg *messaging.Message) *core.Status {
			status := msg.Status()
			if status == nil {
				return core.NewStatusError(core.StatusInvalidArgument, errors.New("message body content is not of type *core.Status"))
			}
			ctx, cancel := context.WithTimeout(context.Background(), addStatusDuration)
			defer cancel()
			return assignment1.InsertStatus(ctx, msg.From(), core.Origin{
				Region:  msg.Header.Get(core.RegionKey),
				Zone:    msg.Header.Get(core.ZoneKey),
				SubZone: msg.Header.Get(core.SubZoneKey),
				Host:    msg.Header.Get(core.HostKey),
			}, status)
		},
	}
}
