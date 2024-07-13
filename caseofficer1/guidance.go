package caseofficer1

import (
	"time"
)

const (
	policyDuration = time.Second * 5
)

type guidance struct {
	policyInterval func() time.Duration
}

func newGuidance() *guidance {
	return &guidance{
		policyInterval: func() time.Duration { return policyDuration },
	}
}
