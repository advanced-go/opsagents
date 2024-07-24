package logistics1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	fmt2 "github.com/advanced-go/stdlib/fmt"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

type testAgent struct{}

func newTestAgent() *testAgent                    { return new(testAgent) }
func (t *testAgent) Uri() string                  { return "testAgent" }
func (t *testAgent) Message(m *messaging.Message) { fmt.Printf("test: testAgent.Message() -> %v\n", m) }
func (t *testAgent) Handle(status *core.Status, _ string) *core.Status {
	fmt.Printf("test: opsAgent.Handle() -> [status:%v]\n", status)
	status.Handled = true
	return status
}
func (t *testAgent) Run()      {}
func (t *testAgent) Shutdown() {}

func newTestOperations() *operations {
	return &operations{
		log: func(handler messaging.OpsAgent, content any) {
			fmt.Printf("test: activity1.Log() -> %v : %v : %v\n", fmt2.FmtRFC3339Millis(time.Now().UTC()), handler.Uri(), content)
		},
	}
}
func testLog(agentId string, content any) *core.Status {
	fmt.Printf("test: activity1.Log() -> %v : %v : %v\n", fmt2.FmtRFC3339Millis(time.Now().UTC()), agentId, content)
	return core.StatusOK()
}

func ExampleAgentUri() {
	u := AgentUri("west")
	fmt.Printf("test: AgentUri() -> [%v]\n", u)

	//Output:
	//test: AgentUri() -> [logistics1:west]

}

func ExampleNewAgent() {
	// a := NewAgent()
	fmt.Printf("test: newAgent() -> ")

	//Output:
	//test: newAgent() ->

}
