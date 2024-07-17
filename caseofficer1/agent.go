package caseofficer1

import (
	"fmt"
	"github.com/advanced-go/intelagents/guidance1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	Class = "case-officer1"
)

type caseOfficer struct {
	running     bool
	uri         string
	ticker      *messaging.Ticker
	traffic     string
	origin      core.Origin
	ctrlC       chan *messaging.Message
	statusCtrlC chan *messaging.Message
	statusC     chan *messaging.Message
	dataCtrlC   chan *messaging.Message
	dataC       chan *messaging.Message
	handler     messaging.OpsAgent
	controllers *messaging.Exchange
	policy      messaging.Agent
	shutdown    func()
}

func AgentUri(traffic string, origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", Class, traffic, origin.Region, origin.Zone)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", Class, traffic, origin.Region, origin.Zone, origin.SubZone)
}

// NewAgent - create a new case officer agent
func NewAgent(interval time.Duration, traffic string, origin core.Origin, handler messaging.OpsAgent) messaging.OpsAgent {
	return newAgent(interval, traffic, origin, handler)
}

// newCAgent - create a new case officer agent
func newAgent(interval time.Duration, traffic string, origin core.Origin, handler messaging.OpsAgent) *caseOfficer {
	c := new(caseOfficer)
	c.uri = AgentUri(traffic, origin)
	c.traffic = traffic
	c.origin = origin
	c.ticker = messaging.NewTicker(interval)

	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.statusCtrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.statusC = make(chan *messaging.Message, 3*messaging.ChannelSize)
	c.dataCtrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.dataC = make(chan *messaging.Message, 3*messaging.ChannelSize)

	c.policy = guidance1.NewPolicyAgent(newGuidance().policyInterval(), c)
	c.handler = handler
	c.controllers = messaging.NewExchange()
	return c
}

// String - identity
func (a *caseOfficer) String() string {
	return a.uri
}

// Uri - agent identifier
func (a *caseOfficer) Uri() string {
	return a.uri
}

// Message - message the agent
func (a *caseOfficer) Message(m *messaging.Message) {
	messaging.Mux(m, a.ctrlC, a.dataC, a.statusC)
}

// Handle - error handler
func (a *caseOfficer) Handle(status *core.Status, requestId string) *core.Status {
	// TODO : do we need any processing specific to a case officer? If not then forward to handler
	return a.handler.Handle(status, requestId)
}

// Shutdown - shutdown the agent
func (a *caseOfficer) Shutdown() {
	if !a.running {
		return
	}
	a.running = false
	if a.shutdown != nil {
		a.shutdown()
	}
	msg := messaging.NewControlMessage(a.uri, a.uri, messaging.ShutdownEvent)
	if a.ctrlC != nil {
		a.ctrlC <- msg
	}
	if a.statusCtrlC != nil {
		a.statusCtrlC <- msg
	}
	a.controllers.Broadcast(msg)
}

// Run - run the agent
func (a *caseOfficer) Run() {
	if a.running {
		return
	}
	a.running = true
	go runStatus(a, activityLog, newAssignment())
	go run(a, activityLog, newControllerAgent, newAssignment())
}

/*
func (c *caseOfficer) startTicker(interval time.Duration) {
	if interval <= 0 {
		interval = c.interval
	} else {
		c.interval = interval
	}
	if c.ticker != nil {
		c.ticker.Stop()
	}
	c.ticker = time.NewTicker(interval)
}

func (c *caseOfficer) stopTicker() {
	c.ticker.Stop()
}


*/
