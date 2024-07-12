package caseofficer1

import (
	"errors"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

// run - case officer data processing
func runData(c *caseOfficer, log logFunc) {
	if c == nil || log == nil {
		return
	}
	for {
		select {
		case msg := <-c.dataC:
			log(c.uri, "processing data message")
			uri := msg.ForwardTo()
			if uri == "" {
				c.Handle(core.NewStatusError(core.StatusInvalidArgument, errors.New("error: invalid message destination")), "")
				continue
			}
			msg.Header.Set(messaging.XTo, uri)
			msg.Header.Set(messaging.XFrom, c.uri)
			err := c.controllers.Send(msg)
			if err != nil {
				c.Handle(core.NewStatusError(core.StatusInvalidArgument, err), "")
			}
		case msg1 := <-c.dataCtrlC:
			switch msg1.Event() {
			case messaging.ShutdownEvent:
				log(c.uri, messaging.ShutdownEvent)
				close(c.dataC)
				close(c.dataCtrlC)
				return
			default:
			}
		default:
		}
	}
}
