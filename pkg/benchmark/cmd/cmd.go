package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var NumWorkers int
var MigrationsPath string
var DataPath string

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Benchmark",
	Long:  "Benchmark",
	Run:   benchmark,
}

func init() {
	benchmarkCmd.PersistentFlags().IntVarP(&NumWorkers, "number-of-workers", "n", 4, "Number of concurrent workers for processing queries")
	benchmarkCmd.PersistentFlags().StringVarP(&MigrationsPath, "migrations", "m", "", "If present, file from which to run migrations. Run no migrations if omitted")
	benchmarkCmd.PersistentFlags().StringVarP(&DataPath, "data", "d", "", "If present, file from which to import data. Run no migrations if omitted")

}

func Execute() {
	if err := benchmarkCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
