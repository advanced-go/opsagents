package logistics1

import (
	"fmt"
	"github.com/advanced-go/operations/activity1"
	"github.com/advanced-go/opsagents/guidance"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	Class = "logistics1"
)

type logistics struct {
	running             bool
	uri                 string
	region              string
	interval            time.Duration
	ticker              *time.Ticker
	caseOfficerInterval time.Duration
	ctrlC               chan *messaging.Message
	caseOfficers        *messaging.Exchange
	shutdown            func()
}

func AgentUri(region string) string {
	return fmt.Sprintf("%v:%v", Class, region)
}

// NewAgent - create a new logistics agent, region needs to be set via host environment
func NewAgent(region string) messaging.OpsAgent {
	return newAgent(region)
}

// newAgent - create a new logistics agent
func newAgent(region string) *logistics {
	c := new(logistics)
	c.uri = AgentUri(region)
	c.region = region
	c.interval = guidance.LogisticsInterval()
	c.caseOfficerInterval = guidance.CaseOfficerInterval()
	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.caseOfficers = messaging.NewExchange()
	return c
}

// String - identity
func (l *logistics) String() string {
	return l.uri
}

// Uri - agent identifier
func (l *logistics) Uri() string {
	return l.uri
}

// Message - message the agent
func (l *logistics) Message(m *messaging.Message) {
	messaging.Mux(m, l.ctrlC, nil, nil)
}

// Handle - error handler
func (c *logistics) Handle(status *core.Status, requestId string) *core.Status {
	// TODO : how to handle an error
	fmt.Printf("test: opsAgent.Handle() -> [status:%v]\n", status)
	status.Handled = true
	return status
}

// Shutdown - shutdown the agent
func (l *logistics) Shutdown() {
	if !l.running {
		return
	}
	l.running = false
	if l.shutdown != nil {
		l.shutdown()
	}
	msg := messaging.NewControlMessage(l.uri, l.uri, messaging.ShutdownEvent)
	if l.ctrlC != nil {
		l.ctrlC <- msg
	}
	l.caseOfficers.Broadcast(msg)
}

// Run - run the agent
func (l *logistics) Run() {
	if l.running {
		return
	}
	l.running = true

	go run(l, activity1.Log, queryAssignments, newCaseOfficer)
}

func (l *logistics) startTicker(interval time.Duration) {
	if interval <= 0 {
		interval = l.interval
	} else {
		l.interval = interval
	}
	if l.ticker != nil {
		l.ticker.Stop()
	}
	l.ticker = time.NewTicker(interval)
}

func (l *logistics) stopTicker() {
	l.ticker.Stop()
}
