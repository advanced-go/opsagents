package logistics1

import (
	"github.com/advanced-go/operations/caseofficer1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

// run - operations logistics
func runLogistics(l *logistics, ls *landscape, ops *operations) {
	if l == nil || ls == nil || ops == nil {
		return
	}
	processAssignments(l, ls)
	ops.log(l, "process assignments : init")
	l.ticker.Start(0)
	for {
		select {
		case <-l.ticker.C():
			// TODO : determine how to check for partition changes
			ops.log(l, "process assignments : tick")
		case msg, open := <-l.ctrlC:
			if !open {
				return
			}
			switch msg.Event() {
			case messaging.ShutdownEvent:
				l.shutdown()
				ops.log(l, messaging.ShutdownEvent)
				return
			default:
			}
		default:
		}
	}
}

func processAssignments(l *logistics, ls *landscape) *core.Status {
	entries, status := ls.query(l.region)
	if !status.OK() {
		l.Handle(status, "")
		return status
	}
	for _, e1 := range entries {
		err := l.caseOfficers.Register(caseofficer1.NewAgent(l.caseOfficerInterval, e1.Traffic, e1.Origin(), l))
		if err != nil {
			status = core.NewStatusError(core.StatusInvalidArgument, err)
			l.Handle(status, "")
			return status
		}
	}
	return core.StatusOK()
}
