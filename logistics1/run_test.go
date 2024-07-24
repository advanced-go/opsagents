package logistics1

import (
	"fmt"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

func ExampleRun() {
	msg := messaging.NewControlMessage("to", "from", messaging.ShutdownEvent)

	c := newAgent("west")
	go runLogistics(c, newLandscape(), newTestOperations())
	time.Sleep(time.Second * 3)

	c.ctrlC <- msg
	time.Sleep(time.Second * 2)

	fmt.Printf("test: run() -> [agents:%v]\n", c.caseOfficers.List())

	//Output:
	//test: activity1.Log() -> 2024-07-08T17:11:39.367Z : logistics1:west : process assignments : init
	//test: activity1.Log() -> 2024-07-08T17:11:41.373Z : logistics1:west : process assignments : tick
	//test: activity1.Log() -> 2024-07-08T17:11:42.372Z : logistics1:west : event:shutdown
	//test: run() -> [agents:[case-officer1:egress.us-central1.c case-officer1:egress.us-central1.d case-officer1:ingress.us-central1.c case-officer1:ingress.us-central1.d]]

}
