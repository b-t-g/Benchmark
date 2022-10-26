package query

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/b-t-g/benchmark/pkg/benchmark/statistics"
	"github.com/jackc/pgx/v4/pgxpool"
)

func RunQuery(queries []string, pool *pgxpool.Pool) statistics.Statistics {
	stats := statistics.Statistics{}

	for _, query := range queries {
		ctx, cFunc := context.WithTimeout(context.Background(), 10*time.Second)
		defer cFunc()

		conn, err := pool.Acquire(ctx)
		if err != nil {
			log.Fatalf("Error acquiring connection from pool: %v", err)
		}

		queryStart := time.Now()
		_, err = conn.Query(ctx, query)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
			os.Exit(1)
		}

		queryEnd := time.Now()
		queryDuration := queryEnd.Sub(queryStart).Milliseconds()
		stats.SumOfDurations += queryDuration

		conn.Release()

		stats.Durations = append(stats.Durations, queryDuration)
		if queryDuration < stats.Min.QueryMsDuration || stats.Min.QueryMsDuration == 0 {
			stats.Min = statistics.QueryStatistic{QueryMsDuration: queryDuration, Query: query}
		}

		if queryDuration > stats.Max.QueryMsDuration {
			stats.Max = statistics.QueryStatistic{QueryMsDuration: queryDuration, Query: query}
		}

	}
	return stats
}
