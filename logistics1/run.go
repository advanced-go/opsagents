package logistics1

import (
	"github.com/advanced-go/operations/caseofficer1"
	"github.com/advanced-go/operations/landscape1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

type logFunc func(agentId string, content any) *core.Status
type agentFunc func(interval time.Duration, traffic string, origin core.Origin, handler messaging.Agent) messaging.Agent
type queryFunc func(region string) ([]landscape1.Entry, *core.Status)

// run - operations logistics
func run(l *logistics, log logFunc, agent agentFunc, ls *landscape) {
	if l == nil || log == nil || agent == nil || ls == nil {
		return
	}
	status := processAssignments(l, agent, ls)
	log(l.uri, "process assignments : init")
	if !status.OK() && !status.NotFound() {
		l.Handle(status, "")
	}
	l.ticker.Start(0)
	for {
		select {
		case <-l.ticker.C():
			// TODO : determine how to check for partition changes
			log(l.uri, "process assignments : tick")
		case msg, open := <-l.ctrlC:
			if !open {
				return
			}
			switch msg.Event() {
			case messaging.ShutdownEvent:
				close(l.ctrlC)
				l.ticker.Stop()
				log(l.uri, messaging.ShutdownEvent)
				return
			default:
			}
		default:
		}
	}
}

func newCaseOfficer(interval time.Duration, traffic string, origin core.Origin, handler messaging.Agent) messaging.Agent {
	return caseofficer1.NewAgent(interval, traffic, origin, handler)
}

func processAssignments(l *logistics, newAgent agentFunc, ls *landscape) *core.Status {
	entries, status := ls.query(l.region)
	if !status.OK() {
		return status
	}
	for _, e1 := range entries {
		err := l.caseOfficers.Register(newAgent(l.caseOfficerInterval, e1.Traffic, e1.Origin(), l))
		if err != nil {
			return core.NewStatusError(core.StatusInvalidArgument, err)
		}
	}
	return status
}
