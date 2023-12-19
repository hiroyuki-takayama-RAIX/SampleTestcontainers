package sample

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	db *sql.DB
)

const (
	driver = "pgx"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithInitScripts("../testdata/init-user-db.sh"),
		postgres.WithConfigFile("../testdata/my-postgres.conf"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForExposedPort(),
			/*
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(5*time.Second)),
			*/
		))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	connStr, err := postgresContainer.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open(driver, connStr)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestQueryRowContext(t *testing.T) {
	q := `SELECT * FROM testdb;`
	row := db.QueryRowContext(context.Background(), q)

	var (
		id   int
		name string
	)
	err := row.Scan(&id, &name)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 1, id)
	assert.Equal(t, "test", name)
}
