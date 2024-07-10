package caseofficer1

import (
	"github.com/advanced-go/agency/egress1"
	"github.com/advanced-go/agency/ingress1"
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

type logFunc func(agentId string, content any) *core.Status
type agentFunc func(traffic string, origin core.Origin, opsAgent messaging.OpsAgent) messaging.Agent

// run - case officer
func run(c *caseOfficer, log logFunc, agent agentFunc, assign *assignment) {
	if c == nil || log == nil || agent == nil || assign == nil {
		return
	}
	status := processAssignments(c, agent, assign)
	log(c.uri, "process assignments : init")
	if !status.OK() && !status.NotFound() {
		c.opsAgent.Handle(status, c.uri)
	}
	c.startTicker(0)
	for {
		select {
		case <-c.ticker.C:
			status = processAssignments(c, agent, assign)
			log(c.uri, "process assignments : tick")
			if !status.OK() && !status.NotFound() {
				c.opsAgent.Handle(status, c.uri)
			}
		case msg, open := <-c.ctrlC:
			if !open {
				return
			}
			switch msg.Event() {
			case messaging.ShutdownEvent:
				close(c.ctrlC)
				c.stopTicker()
				log(c.uri, messaging.ShutdownEvent)
				return
			default:
			}
		default:
		}
	}
}

func newControllerAgent(traffic string, origin core.Origin, opsAgent messaging.OpsAgent) messaging.Agent {
	if traffic == access.IngressTraffic {
		return ingress1.NewControllerAgent(origin, opsAgent)
	}
	return egress1.NewControllerAgent(origin, opsAgent)
}

func processAssignments(c *caseOfficer, agent agentFunc, assign *assignment) *core.Status {
	entries, status := assign.update(c.uri, c.origin)
	if !status.OK() {
		return status
	}
	for _, e := range entries {
		err := c.controllers.Register(agent(c.traffic, e.Origin(), c.opsAgent))
		if err != nil {
			return core.NewStatusError(core.StatusInvalidArgument, err)
		}
	}
	return status
}
