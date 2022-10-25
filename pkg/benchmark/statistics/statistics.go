package statistics

import (
	"math"
	"sort"
)

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

func ProcessQueryStatistics(stats Statistics) ProcessedStatistics {
	r := ProcessedStatistics{}
	r.Median = median(stats.Durations)
	r.Average = average(stats.Durations)
	r.StdDev = standardDeviation(stats.Durations, r.Average)
	r.Min = stats.Min
	r.Max = stats.Max
	return r
}

func standardDeviation(durations []int64, avg float64) float64 {
	var sd float64
	for _, i := range durations {
		sd += math.Pow(float64(i) - avg, 2)
	}

	return math.Sqrt(sd/(float64(len(durations))))
}

func average(durations []int64) float64 {
	avg := int64(0)
	for _, i := range durations {
		avg += i
	}

	return float64(avg)/float64(len(durations))
}

func median(durations []int64) float64 {
	sort.Slice(durations, func(i,j int) bool {
		return durations[i] < durations[j]
	})

	if len(durations) % 2 == 1 {
		return float64(durations[(len(durations) - 1)/2])
	}

	return (float64(durations[(len(durations)/2) - 1]) + float64(durations[(len(durations)/2)]))/2
}

