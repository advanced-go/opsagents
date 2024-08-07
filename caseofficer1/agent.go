package caseofficer1

import (
	"fmt"
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

	//c.policy = guidance1.NewPolicyAgent(newGuidance().policyInterval(), c)
	c.handler = handler
	c.controllers = messaging.NewExchange()
	return c
}

// String - identity
func (c *caseOfficer) String() string {
	return c.uri
}

// Uri - agent identifier
func (c *caseOfficer) Uri() string {
	return c.uri
}

// Message - message the agent
func (c *caseOfficer) Message(m *messaging.Message) {
	messaging.Mux(m, c.ctrlC, c.dataC, c.statusC)
}

// Handle - error handler
func (c *caseOfficer) Handle(status *core.Status, requestId string) *core.Status {
	// TODO : do we need any processing specific to a case officer? If not then forward to handler
	return c.handler.Handle(status, requestId)
}

// AddActivity - add activity
func (c *caseOfficer) AddActivity(agentId string, content any) {
	// TODO : Any operations specific processing ??  If not then forward to handler
	//return a.handler.Handle(status, requestId)
}

// Shutdown - shutdown the agent
func (c *caseOfficer) Shutdown() {
	if !c.running {
		return
	}
	c.running = false
	if c.shutdown != nil {
		c.shutdown()
	}
	msg := messaging.NewControlMessage(c.uri, c.uri, messaging.ShutdownEvent)
	if c.ctrlC != nil {
		c.ctrlC <- msg
	}
	if c.statusCtrlC != nil {
		c.statusCtrlC <- msg
	}
	c.controllers.Broadcast(msg)
}

// Run - run the agent
func (c *caseOfficer) Run() {
	if c.running {
		return
	}
	c.running = true
	go runStatus(c, activityLog, newAssignment())
	go runCaseOfficer(c, activityLog, newControllerAgent, newAssignment())
}
