package guidance

import "time"

func LogisticsInterval() time.Duration {
	return time.Second * 2
}

func CaseOfficerInterval() time.Duration {
	return time.Second * 5
}
