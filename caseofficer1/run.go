package caseofficer1

import (
	"context"
	"github.com/advanced-go/agency/egress1"
	"github.com/advanced-go/agency/ingress1"
	"github.com/advanced-go/operations/assignment1"
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

type logFunc func(ctx context.Context, agentId string, content any) *core.Status
type updateFunc func(ctx context.Context, agentId string, origin core.Origin) ([]assignment1.Entry, *core.Status)
type agentFunc func(traffic string, origin core.Origin, opsAgent messaging.OpsAgent) messaging.Agent

// run - case officer
func run(c *caseOfficer, log logFunc, update updateFunc, agent agentFunc) {
	if c == nil {
		return
	}
	status := processAssignments(c, update, agent)
	log(nil, c.uri, "process assignments : init")
	if !status.OK() && !status.NotFound() {
		c.opsAgent.Handle(status, c.uri)
	}
	c.startTicker(0)
	for {
		select {
		case <-c.ticker.C:
			status = processAssignments(c, update, agent)
			log(nil, c.uri, "process assignments : tick")
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
				log(nil, c.uri, messaging.ShutdownEvent)
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

func processAssignments(c *caseOfficer, update updateFunc, agent agentFunc) *core.Status {
	entries, status := update(nil, c.uri, c.origin)
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
