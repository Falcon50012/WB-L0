package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

var DBPool *pgxpool.Pool

func InitDBConnectionPool() {
	connectionData := "postgres://wbdev:123@localhost:5432/postgres"
	DB, err := pgxpool.New(context.Background(), connectionData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	} else {
		log.Println("Connection pool created")
	}
	DBPool = DB
}

func CreateTables() error {
	path := "database/migration/structure.sql"

	c, ioErr := os.ReadFile(path)
	if ioErr != nil {
		return fmt.Errorf("ERROR: can't find init sql file: %v", ioErr)
	}

	conn, err := DBPool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("ERROR: can't return connection from the Pool: %v", err)
	}
	defer conn.Release()

	sql := string(c)
	_, execErr := conn.Exec(context.Background(), sql)
	if execErr != nil {
		return fmt.Errorf("ERROR: can't run init script: %v", execErr)
	}

	log.Println("Tables initialized successfully")
	return nil
}
