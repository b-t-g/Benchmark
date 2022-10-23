package statistics

import "time"

type Statistics struct {
	min       time.Duration
	max       time.Duration
	durations []time.Duration
}
