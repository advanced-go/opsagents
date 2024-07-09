package caseofficer1

import (
	"errors"
	"github.com/advanced-go/operations/assignment1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

type insertFunc func(msg *messaging.Message) *core.Status

// run - case officer status processing
func runStatus(c *caseOfficer, log logFunc, insert insertFunc) {
	if c == nil {
		return
	}
	for {
		select {
		case msg, open := <-c.statusC:
			if !open {
				return
			}
			log(nil, c.uri, "processing status message")
			status := insert(msg)
			if !status.OK() && !status.NotFound() {
				c.handler.Message(messaging.NewStatusMessage(c.handler.Uri(), c.uri, status))
			}
		case msg1, open1 := <-c.statusCtrlC:
			if !open1 {
				return
			}
			switch msg1.Event() {
			case messaging.ShutdownEvent:
				log(nil, c.uri, messaging.ShutdownEvent)
				close(c.statusC)
				close(c.statusCtrlC)
				return
			default:
			}
		default:
		}
	}
}

func insertAssignmentStatus(msg *messaging.Message) *core.Status {
	status := msg.Status()
	if status == nil {
		return core.NewStatusError(core.StatusInvalidArgument, errors.New("message body content is not of type *core.Status"))
	}
	return assignment1.InsertStatus(nil, msg.From(), core.Origin{
		Region:  msg.Header.Get(core.RegionKey),
		Zone:    msg.Header.Get(core.ZoneKey),
		SubZone: msg.Header.Get(core.SubZoneKey),
		Host:    msg.Header.Get(core.HostKey),
	}, status)
}
