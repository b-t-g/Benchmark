package statistics

import "time"

type Statistics struct {
	Min       QueryStatistic
	Max       QueryStatistic
	Durations []int64
}

type QueryStatistic struct {
	QueryMsDuration int64
	Query           string
}
