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
		It("fails to validate a row where start date is after end date", func() {
			err := cmd.ValidateRow("host_000008,2017-01-01 09:59:22,2017-01-01 08:59:22")
			Expect(err).To(HaveOccurred())
		})
		It("fails to validate a row with more than 3 columns", func() {
			err := cmd.ValidateRow("host name with,comma,2017-01-01 08:59:22,2017-01-01 09:59:22")
			Expect(err).To(HaveOccurred())
		})
		It("fails to validate a row with fewer than 3 columns", func() {
			err := cmd.ValidateRow("2017-01-01 08:59:22,2017-01-01 09:59:22")
			Expect(err).To(HaveOccurred())
		})
		It("fails to validate a with the columns in the wrong order", func() {
			err := cmd.ValidateRow("2017-01-01 08:59:22, host_000008,2017-01-01 09:59:22")
			Expect(err).To(HaveOccurred())
		})
	})
})
