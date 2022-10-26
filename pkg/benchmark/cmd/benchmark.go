package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/b-t-g/benchmark/pkg/benchmark/query"
	"github.com/b-t-g/benchmark/pkg/benchmark/statistics"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"
)

const (
	username = "postgres"
	password = "example"
	port     = 5432
	database = "homework"
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
	threadsToQueries := make([][]string, NumWorkers)
	for i := 0; i < NumWorkers; i++ {
		threadsToQueries[i] = []string{}
	}

	hostToThread := map[string]int{}

	// Assign queries to threads in a round robin fashion
	threadScheduler := 0
	totalQueries := 0

	// Skip the first line defining columns
	scanner.Scan()

	for scanner.Scan() {
		text = scanner.Text()
		if err = ValidateRow(text); err != nil {
			fmt.Printf("Skipping row %s for reason %v\n", text, err)
			continue
		}
		totalQueries += 1

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

		threadsToQueries[thread] = append(threadsToQueries[thread], formatQueryString(host, startTime, endTime))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hostname := os.Getenv("DB_HOSTNAME")
	if hostname == "" {
		log.Fatal("DB_HOSTNAME for declaring the DB hostname is empty!")
	}

	// In a Kubernetes environment, instead of hard-coding, I'd create a Kubernetes
	// secret and, in both the stateful set for the database (if it's deployed myself)
	// and the deployment for this benchmark tool, I'd use the secret as an environment
	// variable as shown here:
	// https://kubernetes.io/docs/concepts/configuration/secret/#using-secrets-as-environment-variables
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		username, password, hostname, port, database, sslmode)
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	wg := new(sync.WaitGroup)
	queryStatistics := statistics.Statistics{}
	queryStatisticsMutex := sync.Mutex{}
	var processStartTime time.Time
	for i := 0; i < NumWorkers; i++ {
		wg.Add(1)
		var localQueryStatistics statistics.Statistics

		// asynchronous functions work oddly with loop variables
		j := i
		processStartTime = time.Now()
		go func() {
			localQueryStatistics = query.RunQuery(threadsToQueries[j], pool)

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
	processEndTime := time.Now()
	processingTime := processEndTime.Sub(processStartTime).Milliseconds()

	processedStatistics := statistics.ProcessQueryStatistics(queryStatistics)
	fmt.Printf("%s", formatOutput(processedStatistics, processingTime, totalQueries))
	os.Exit(0)
}

func formatOutput(processedStatistics statistics.ProcessedStatistics, processingTime int64, totalQueries int) string {
	return fmt.Sprintf(`
Total Processing Time: %d ms
Total Queries Processed: %d

Min Query Time: %d ms
Query with Min Time: %s

Max Query Time: %d ms
Query With Max Time: %s

Average Query Time: %f ms

Median Query Time: %f ms

Standard Deviation in Query Time: %f ms
`, processingTime,
		totalQueries,
		processedStatistics.Min.QueryMsDuration,
		processedStatistics.Min.Query,
		processedStatistics.Max.QueryMsDuration,
		processedStatistics.Max.Query,
		processedStatistics.Average,
		processedStatistics.Median,
		processedStatistics.StdDev)
}

func ValidateRow(row string) error {
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
