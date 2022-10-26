package cmd_test

import (
	"github.com/b-t-g/benchmark/pkg/benchmark/cmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Benchmark", func() {
	Context("When validating rows", func() {
		It("validates a valid row", func() {
			err := cmd.ValidateRow("host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22")
			Expect(err).To(BeNil())
		})
		It("fails to validate a row with both dates invalid", func() {
			err := cmd.ValidateRow("host_000008,2017-01-1 08:59:22,2017-01-1 09:59:22")
			Expect(err).To(HaveOccurred())
		})
		It("fails to validate a row with one date invalid", func() {
			err := cmd.ValidateRow("host_000008,2017-01-1 08:59:22,2017-01-01 09:59:22")
			Expect(err).To(HaveOccurred())
		})
		It("fails to validate a row with one date that looks nothing like a date", func() {
			err := cmd.ValidateRow("host_000008,2017-01-1 08:59:22,hello!")
			Expect(err).To(HaveOccurred())
		})
	})
})
