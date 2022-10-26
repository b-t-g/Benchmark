package statistics_test

import (
	"math"

	"github.com/b-t-g/benchmark/pkg/benchmark/statistics"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Statistics", func() {
	var stats statistics.Statistics

	Context("When calculating statistics", func() {
		BeforeEach(func() {
			stats.Min = statistics.QueryStatistic{QueryMsDuration: 1}
			stats.Max = statistics.QueryStatistic{QueryMsDuration: 5}
			stats.Durations = []int64{1, 2, 3, 4, 5}
		})
		It("calculates statistics correctly", func() {
			p := statistics.ProcessQueryStatistics(stats)
			Expect(p.Average).To(Equal(float64(3)))
			Expect(p.Median).To(Equal(float64(3)))
			Expect(p.Min.QueryMsDuration).To(Equal(int64(1)))
			Expect(p.Max.QueryMsDuration).To(Equal(int64(5)))
			Expect(p.StdDev).To(Equal(math.Sqrt(2)))
		})
		It("calculates statistics the same even when durations are unsorted", func() {
			p1 := statistics.ProcessQueryStatistics(stats)

			stats.Durations = []int64{1, 4, 3, 5, 2}
			p2 := statistics.ProcessQueryStatistics(stats)
			Expect(p1.Average).To(Equal(p2.Average))
			Expect(p1.Median).To(Equal(p2.Median))
			Expect(p1.Min.QueryMsDuration).To(Equal(p2.Min.QueryMsDuration))
			Expect(p1.Max.QueryMsDuration).To(Equal(p2.Max.QueryMsDuration))
			Expect(p1.StdDev).To(Equal(p2.StdDev))
		})
	})
})
