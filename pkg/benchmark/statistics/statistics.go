package statistics

import (
	"math"
	"sort"
)

type Statistics struct {
	Min            QueryStatistic
	Max            QueryStatistic
	SumOfDurations int64
	Durations      []int64
}

type QueryStatistic struct {
	QueryMsDuration int64

	// Include the whole query, instead of just the query params
	// so users can compare results just by copy/pasting the query
	// and run it themselves, instead of needing to construct the query
	// manually from the parameters.
	Query string
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
	r.Average = average(stats.SumOfDurations, len(stats.Durations))
	r.StdDev = standardDeviation(stats.Durations, r.Average)
	r.Min = stats.Min
	r.Max = stats.Max
	return r
}

// Include standard deviation to be able to better interpret the other statistics,
// especially min and max. E.g., determine whether min and max are outliers or if
// most of the queries are just very spread out.
func standardDeviation(durations []int64, avg float64) float64 {
	var sd float64
	for _, i := range durations {
		sd += math.Pow(float64(i)-avg, 2)
	}

	return math.Sqrt(sd / (float64(len(durations))))
}

func average(sumOfDurations int64, lenDurations int) float64 {
	return float64(sumOfDurations) / float64(lenDurations)
}

// This is the simplest approach, and hence easiest to understand/change, but, in terms of time complexity, is not optimal.
// If improving statistics calculation time could lead to an appreciably better experience, we have a few options:
//
//  1. Use a better algorithm, namely the "quickselect" algorithm with the (expected) complexity
//     of O(N) where N is the length of the array.
//
//  2. Keep the algorithms for calculating statistics largely the same, but amortize their calculation. I.e.,
//     when sorting the array (e.g., using heapsort, which analyzes each element individually by "removing" each element
//     from the heap), we could calculate the average at the same time. This would not lead to a runtime improvement in terms
//     of computational complexity, but could be an appreciable improvement to a user when there are a sufficiently large number
//     of queries.
//
// Given that I'd expect the runtime to be dominated by the network latency of the queries in most cases, I felt that
// a simpler to understand approach for computing the median, though with worse runtime complexity, was the correct tradeoff.
func median(durations []int64) float64 {
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	if len(durations)%2 == 1 {
		return float64(durations[(len(durations)-1)/2])
	}

	return (float64(durations[(len(durations)/2)-1]) + float64(durations[(len(durations)/2)])) / 2
}
