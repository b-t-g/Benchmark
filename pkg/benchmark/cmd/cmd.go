package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var NumWorkers int
var QueryParamsPath string

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Benchmark",
	Long:  "Benchmark",
	Run:   benchmark,
}

func init() {
	benchmarkCmd.PersistentFlags().IntVarP(&NumWorkers, "number-of-workers", "n", 4, "Number of concurrent workers for processing queries")
	benchmarkCmd.PersistentFlags().StringVarP(&QueryParamsPath, "query-params", "q", "", "Required: File from which to read to generate queries.")

}

func Execute() {
	if err := benchmarkCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
