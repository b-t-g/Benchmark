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
	for scanner.Scan() {
		// Just get the last host for now
		text = scanner.Text()
		if strings.Split(text, ",")[0] == "hostname" {
			continue
		}
		validateRow(text)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// determined to be valid
	host := strings.Split(text, ",")[0]
	startTime := strings.Split(text, ",")[1]
	endTime := strings.Split(text, ",")[2]
	queryString := formatQueryString(host, startTime, endTime)

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
	fmt.Println("Done")
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
