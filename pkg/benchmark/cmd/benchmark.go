package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/b-t-g/benchmark/pkg/benchmark/query"
	"github.com/spf13/cobra"
)

const (
	username       = "postgres"
	password       = "example"
	port           = 5432
	database       = "postgres"
	sslmode        = "disable"
	sampleQueryFmt = `
select time_bucket('1 minute', ts) as one_min, min(usage), max(usage) from cpu_usage
where host = '%s'
group by one_min   
order by one_min desc limit 5;
`
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

	scanner := bufio.NewScanner(file)

	var text string
	for scanner.Scan() {
		// Just get the last host for now
		text = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	host := strings.Split(text, ",")[0]
	queryString := fmt.Sprintf(sampleQueryFmt, host)

	wg := new(sync.WaitGroup)
	for i := 0; i < NumWorkers; i++ {
		wg.Add(1)
		j := i
		go func() {
			query.RunQuery(j, queryString)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("Done")
}
