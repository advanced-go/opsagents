package caseofficer1

import (
	"github.com/advanced-go/intelagents/egress1"
	"github.com/advanced-go/intelagents/ingress1"
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

type logFunc func(agentId string, content any) *core.Status
type agentFunc func(traffic string, origin core.Origin, handler messaging.OpsAgent) messaging.Agent

// run - case officer
func runCaseOfficer(a *caseOfficer, log logFunc, agent agentFunc, assign *assignment) {
	if a == nil || log == nil || agent == nil || assign == nil {
		return
	}
	status := processAssignments(a, agent, assign)
	log(a.uri, "process assignments : init")
	if !status.OK() && !status.NotFound() {
		a.handler.Handle(status, a.uri)
	}
	a.ticker.Start(0)
	for {
		select {
		case <-a.ticker.C():
			status = processAssignments(a, agent, assign)
			log(a.uri, "process assignments : tick")
			if !status.OK() && !status.NotFound() {
				a.handler.Handle(status, a.uri)
			}
		case msg, open := <-a.ctrlC:
			if !open {
				return
			}
			switch msg.Event() {
			case messaging.ShutdownEvent:
				close(a.ctrlC)
				a.ticker.Stop()
				log(a.uri, messaging.ShutdownEvent)
				return
			default:
			}
		default:
		}
	}
}

func newControllerAgent(traffic string, origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	if traffic == access.IngressTraffic {
		return ingress1.NewFieldOperative(origin, nil, handler)
	}
	return egress1.NewLeadAgent(origin, handler)
}

func processAssignments(a *caseOfficer, agent agentFunc, assign *assignment) *core.Status {
	entries, status := assign.update(a.uri, a.origin)
	if !status.OK() {
		return status
	}
	for _, e := range entries {
		err := a.controllers.Register(agent(a.traffic, e.Origin(), a.handler))
		if err != nil {
			return core.NewStatusError(core.StatusInvalidArgument, err)
		}
	}
	return status
}
