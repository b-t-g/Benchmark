package query

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/b-t-g/benchmark/pkg/benchmark/statistics"
	pgx "github.com/jackc/pgx/v4"
)

const (
	username = "postgres"
	password = "example"
	port     = 5432
	database = "homework"
	sslmode  = "disable"
)

func RunQuery(query string) statistics.Statistics {
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
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	stats := statistics.Statistics{}

	queryStart := time.Now()
	_, err = conn.Query(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	queryEnd := time.Now()
	queryDuration := queryEnd.Sub(queryStart).Milliseconds()

	stats.Durations = append(stats.Durations, queryDuration)
	if queryDuration < stats.Min.QueryMsDuration || stats.Min.QueryMsDuration == 0 {
		stats.Min = statistics.QueryStatistic{QueryMsDuration: queryDuration, Query: query}
	}

	if queryDuration > stats.Max.QueryMsDuration {
		stats.Max = statistics.QueryStatistic{QueryMsDuration: queryDuration, Query: query}
	}

	return stats
}
