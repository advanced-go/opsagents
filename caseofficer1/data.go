package caseofficer1

import "github.com/advanced-go/stdlib/messaging"

// run - case officer data processing
func runData(c *caseOfficer, log logFunc) {
	if c == nil || log == nil {
		return
	}
	for {
		select {
		case msg := <-c.dataC:
			log(c.uri, "processing data message")
			if uri := msg.Event(); uri != "" {
			}

		case msg1 := <-c.dataCtrlC:
			switch msg1.Event() {
			case messaging.ShutdownEvent:
				log(c.uri, messaging.ShutdownEvent)
				close(c.dataC)
				close(c.dataCtrlC)
				return
			default:
			}
		default:
		}
	}
}
