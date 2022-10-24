package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/b-t-g/benchmark/pkg/benchmark/query"
	"github.com/spf13/cobra"
)

const (
	username = "postgres"
	password = "example"
	port     = 5432
	database = "postgres"
	sslmode  = "disable"
)

func benchmark(cmd *cobra.Command, args []string) {
	if QueryParamsPath == "" {
		fmt.Println("Missing required field query-params")
		os.Exit(1)
	}
	file, err := os.Open(QueryParamsPath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	wg := new(sync.WaitGroup)
	for i := 0; i < NumWorkers; i++ {
		wg.Add(1)
		j := i
		go func() {
			query.RunQuery(j, "")
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("Done")
}
