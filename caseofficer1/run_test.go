package caseofficer1

import (
	"fmt"
	"github.com/advanced-go/operations/assignment1"
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

func ExampleNewControllerAgent() {
	origin := core.Origin{
		Region:     "us-central1",
		Zone:       "c",
		SubZone:    "",
		Host:       "www.host1.com",
		InstanceId: "",
	}
	a := newControllerAgent(access.IngressTraffic, origin, nil)
	fmt.Printf("test: newControllerAgent(\"%v\") -> [%v]\n", access.IngressTraffic, a)

	a = newControllerAgent(access.EgressTraffic, origin, nil)
	fmt.Printf("test: newControllerAgent(\"%v\") -> [%v]\n", access.EgressTraffic, a)

	//Output:
	//test: newControllerAgent("ingress") -> [ingress-controller1:us-central1.c.www.host1.com]
	//test: newControllerAgent("egress") -> [egress-controller1:us-central1.c.www.host1.com]

}

func ExampleProcessAssignments() {
	origin := core.Origin{
		Region:     "us-central1",
		Zone:       "c",
		SubZone:    "",
		Host:       "www.host1.com",
		InstanceId: "",
	}

	c := newAgent(time.Second*5, access.IngressTraffic, origin, nil)
	fmt.Printf("test: newAgent() -> [status:%v]\n", c != nil)

	status := processAssignments(c, assignment1.Update, newControllerAgent)
	fmt.Printf("test: processAssignments() -> [status:%v] [controllers:%v]\n", status, c.controllers.Count())

	//Output:
	//test: newAgent() -> [status:true]
	//test: processAssignments() -> [status:OK] [controllers:2]

}

func ExampleRun() {
	origin := core.Origin{
		Region:     "us-central1",
		Zone:       "c",
		SubZone:    "",
		Host:       "www.host1.com",
		InstanceId: "",
	}
	msg := messaging.NewControlMessage("to", "from", messaging.ShutdownEvent)

	c := newAgent(time.Second*50, access.IngressTraffic, origin, newTestAgent())
	go run(c, testLog, assignment1.Update, newControllerAgent)
	time.Sleep(time.Second * 2)

	c.ctrlC <- msg
	time.Sleep(time.Second * 3)

	fmt.Printf("test: run() -> [agents:%v]\n", c.controllers.List())

	c.controllers.Broadcast(msg)
	fmt.Printf("test: exchange.Broadcast() -> [agents:%v]\n", c.controllers.List())

	//Output:
	//test: activity1.Log() -> 2024-07-08T15:06:51.286Z : case-officer1:ingress.us-central1.c : process assignments : init
	//test: activity1.Log() -> 2024-07-08T15:06:53.288Z : case-officer1:ingress.us-central1.c : event:shutdown
	//test: run() -> [agents:[ingress-controller1:us-central1.c.www.host1.com ingress-controller1:us-central1.c.www.host2.com]]
	//test: exchange.Broadcast() -> [agents:[]]

}

func ExampleRun_Error() {
	origin := core.Origin{
		Region:     "us-central1",
		Zone:       "c",
		SubZone:    "",
		Host:       "www.host1.com",
		InstanceId: "",
	}
	msg := messaging.NewControlMessage("to", "from", messaging.ShutdownEvent)

	c := newAgent(time.Second*2, access.IngressTraffic, origin, newTestAgent())
	go run(c, testLog, assignment1.Update, newControllerAgent)
	time.Sleep(time.Second * 3)

	c.ctrlC <- msg
	time.Sleep(time.Second * 3)

	//Output:
	//test: activity1.Log() -> 2024-07-08T17:02:46.989Z : case-officer1:ingress.us-central1.c : process assignments : default
	//test: activity1.Log() -> 2024-07-08T17:02:48.993Z : case-officer1:ingress.us-central1.c : process assignments : onTick()
	//test: testAgent.Message() -> [chan:STATUS] [from:case-officer1:ingress.us-central1.c] [to:testAgent] [event:status] [status:Invalid Argument [error: controller2.Register() agent already exists: [ingress-controller1:us-central1.c.www.host1.com]]]
	//test: activity1.Log() -> 2024-07-08T17:02:49.990Z : case-officer1:ingress.us-central1.c : event:shutdown

}
