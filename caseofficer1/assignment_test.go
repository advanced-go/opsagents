package caseofficer1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

func ExampleAssignmentAddStatus() {
	msg := messaging.NewMessageWithStatus(messaging.ChannelStatus, "to", "from", "", core.StatusOK())
	status := newAssignment().addStatus(msg)

	fmt.Printf("test: assignment.AddStatus() -> [status:%v]\n", status)

	//Output:
	//test: insertAssignmentStatus() -> [status:OK]

}
