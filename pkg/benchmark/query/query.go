package query

import (
	"context"
	"fmt"
	"log"
	"os"

	pgx "github.com/jackc/pgx/v4"
)

const (
	username = "postgres"
	password = "example"
	port     = 5432
	database = "postgres"
	sslmode  = "disable"
)

func RunQuery(goRoutineNumber int, query string) {
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
		return
	}
	defer conn.Close(ctx)

	fmt.Printf("%d finished!", goRoutineNumber)
}