package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/b-t-g/benchmark/pkg/benchmark/query"
	"github.com/b-t-g/benchmark/pkg/benchmark/statistics"
	"github.com/spf13/cobra"
)

const (
	username = "postgres"
	password = "example"
	port     = 5432
	database = "postgres"
	sslmode  = "disable"
	queryFmt = `
select time_bucket('1 minute', ts, '%s seconds'::INTERVAL) as one_min, min(usage), max(usage) from cpu_usage
where host = '%s' and ts >= '%s' and ts <= '%s' 
group by one_min ;
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
	threadsToQueries := [][]string{}

	hostToThread := map[string]int

	// Assign queries to threads in a round robin fashion
	threadScheduler := 0

	// Skip the first line defining columns
	scanner.Scan()

	for scanner.Scan() {
		text = scanner.Text()
		if err = validateRow(text); err != nil {
			fmt.Printf("Skipping row %s for reason %v\n", text, err)
			continue
		}

		host := strings.Split(text, ",")[0]
		startTime := strings.Split(text, ",")[1]
		endTime := strings.Split(text, ",")[2]
		var thread int
		if k, found := hostToThread[host]; found {
			thread = k
		} else {
			thread = threadScheduler
			threadScheduler = (threadScheduler + 1) % NumWorkers
		}

		queriesForCurrentThread := threadsToQueries[thread]
		queriesForCurrentThread = append(queriesForCurrentThread, formatQueryString(host, startTime, endTime))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wg := new(sync.WaitGroup)
	queryStatistics := statistics.Statistics{}
	queryStatisticsMutex := sync.Mutex{}
	for i := 0; i < NumWorkers; i++ {
		wg.Add(1)
		var localQueryStatistics statistics.Statistics
		go func() {
			localQueryStatistics = query.RunQuery(queryString)

			queryStatisticsMutex.Lock()
			defer queryStatisticsMutex.Unlock()

			if localQueryStatistics.Min.QueryMsDuration < queryStatistics.Min.QueryMsDuration || queryStatistics.Min.QueryMsDuration == 0 {
				queryStatistics.Min = localQueryStatistics.Min
			}

			if localQueryStatistics.Max.QueryMsDuration > queryStatistics.Max.QueryMsDuration {
				queryStatistics.Max = localQueryStatistics.Max
			}

			queryStatistics.Durations = append(queryStatistics.Durations, localQueryStatistics.Durations...)
			wg.Done()
		}()

	}

	wg.Wait()

	fmt.Printf("Done\n")
	fmt.Printf("%v\n", queryStatistics)
	fmt.Printf("%v\n", statistics.ProcessQueryStatistics(queryStatistics))
	os.Exit(0)
}

func validateRow(row string) error {
	fields := strings.Split(row, ",")
	layout := "2006-02-01 15:04:05"
	_, err := time.Parse(layout, fields[1])
	if err != nil {
		fmt.Printf("Error parsing timestamp %s\nIn row %s\n", fields[1], row)
		return err
	}

	_, err = time.Parse(layout, fields[2])
	if err != nil {
		fmt.Printf("Error parsing timestamp %s\nIn row %s\n", fields[2], row)
		return err
	}

	return nil
}

func formatQueryString(host, startTime, endTime string) string {
	splitStartTime := strings.Split(startTime, ":")
	seconds := splitStartTime[len(splitStartTime)-1]
	return fmt.Sprintf(queryFmt, seconds, host, startTime, endTime)
}
