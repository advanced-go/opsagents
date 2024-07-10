package caseofficer1

import (
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"net/http"
	"time"
)

func ExampleInsertAssignmentStatus() {
	msg := messaging.NewMessageWithStatus(messaging.ChannelStatus, "to", "from", "", core.StatusOK())
	status := insertAssignmentStatus(msg)

	fmt.Printf("test: insertAssignmentStatus() -> [status:%v]\n", status)

	//Output:
	//test: insertAssignmentStatus() -> [status:OK]

}

func ExampleRunStatus() {
	origin := core.Origin{
		Region:     "us-central1",
		Zone:       "c",
		SubZone:    "",
		Host:       "www.host1.com",
		InstanceId: "",
	}
	msg := messaging.NewControlMessage("to", "from", messaging.ShutdownEvent)
	//fmt.Printf("test: NewMessage() -> %v\n", msg.Event())

	c := newAgent(time.Second*1, access.IngressTraffic, origin, newTestAgent())
	go runStatus(c, testLog, insertAssignmentStatus)

	status := core.NewStatusError(http.StatusTeapot, errors.New("teapot error"))
	c.statusC <- messaging.NewMessageWithStatus(messaging.ChannelStatus, "to", "from", "event:status", status)
	time.Sleep(time.Second * 1)

	c.statusCtrlC <- msg
	time.Sleep(time.Second * 3)

	//Output:
	//test: activity1.Log() -> 2024-07-08T14:30:07.973Z : case-officer1:ingress.us-central1.c : processing status message
	//test: activity1.Log() -> 2024-07-08T14:30:08.975Z : case-officer1:ingress.us-central1.c : event:shutdown

}

func ExampleRunStatus_Error() {
	origin := core.Origin{
		Region:     "us-central1",
		Zone:       "c",
		SubZone:    "",
		Host:       "www.host1.com",
		InstanceId: "",
	}
	msg := messaging.NewControlMessage("to", "from", messaging.ShutdownEvent)

	c := newAgent(time.Second*1, access.IngressTraffic, origin, newTestAgent())
	go runStatus(c, testLog, func(m *messaging.Message) *core.Status {
		return core.NewStatusError(http.StatusGatewayTimeout, errors.New("context deadline exceeded"))
	})

	status := core.NewStatusError(http.StatusTeapot, errors.New("teapot error"))
	c.statusC <- messaging.NewMessageWithStatus(messaging.ChannelStatus, "to", "from", "event:status", status)
	time.Sleep(time.Second * 1)

	c.statusCtrlC <- msg
	time.Sleep(time.Second * 3)

	//Output:
	//test: activity1.Log() -> 2024-07-08T14:35:38.921Z : case-officer1:ingress.us-central1.c : processing status message
	//test: opsAgent.Handle() -> [status:Timeout [context deadline exceeded]]
	//test: activity1.Log() -> 2024-07-08T14:35:39.925Z : case-officer1:ingress.us-central1.c : event:shutdown

}
