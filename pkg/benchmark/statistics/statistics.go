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

type ProcessedStatistics struct {
	Min     QueryStatistic
	Max     QueryStatistic
	Median  float64
	Average float64
	StdDev  float64
}
