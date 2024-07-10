package logistics1

import (
	"context"
	"github.com/advanced-go/operations/caseofficer1"
	"github.com/advanced-go/operations/landscape1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"net/url"
	"time"
)

type logFunc func(ctx context.Context, agentId string, content any) *core.Status
type agentFunc func(interval time.Duration, traffic string, origin core.Origin, handler messaging.Agent) messaging.Agent
type queryFunc func(region string) ([]landscape1.Entry, *core.Status)

// run - operations logistics
func run(l *logistics, log logFunc, query queryFunc, agent agentFunc) {
	if l == nil {
		return
	}
	status := processAssignments(l, query, agent)
	log(nil, l.uri, "process assignments : init")
	if !status.OK() && !status.NotFound() {
		l.Handle(status, "")
	}
	l.startTicker(0)
	for {
		select {
		case <-l.ticker.C:
			// TODO : determine how to check for partition changes
			log(nil, l.uri, "process assignments : tick")
		case msg, open := <-l.ctrlC:
			if !open {
				return
			}
			switch msg.Event() {
			case messaging.ShutdownEvent:
				close(l.ctrlC)
				l.stopTicker()
				log(nil, l.uri, messaging.ShutdownEvent)
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

func queryAssignments(region string) ([]landscape1.Entry, *core.Status) {
	values := make(url.Values)
	values.Add(landscape1.AssignedRegionKey, region)
	values.Add(landscape1.StatusKey, landscape1.StatusActive)
	return landscape1.Get(nil, nil, values)
}

func processAssignments(l *logistics, query queryFunc, newAgent agentFunc) *core.Status {
	entries, status := query(l.region)
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
