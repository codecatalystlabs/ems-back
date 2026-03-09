package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	var (
		command       = flag.String("command", "up", "goose command: up, down, status, version, up-by-one, down-to")
		migrationsDir = flag.String("dir", "migrations", "path to migrations directory")
		dbURL         = flag.String("db", envOrDefault("DATABASE_URL", "postgres://dispatch:dispatch@localhost:5431/dispatch_db?sslmode=disable"), "database connection string")
		targetVersion = flag.Int64("version", 0, "target migration version for down-to or up-to")
	)
	flag.Parse()

	db, err := goose.OpenDBWithDriver("pgx", *dbURL)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	goose.SetTableName("goose_db_version")

	switch *command {
	case "up":
		err = goose.Up(db, *migrationsDir)
	case "up-by-one":
		err = goose.UpByOne(db, *migrationsDir)
	case "up-to":
		if *targetVersion == 0 {
			log.Fatal("-version is required for up-to")
		}
		err = goose.UpTo(db, *migrationsDir, *targetVersion)
	case "down":
		err = goose.Down(db, *migrationsDir)
	case "down-to":
		err = goose.DownTo(db, *migrationsDir, *targetVersion)
	case "status":
		err = goose.Status(db, *migrationsDir)
	case "version":
		var v int64
		v, err = goose.GetDBVersion(db)
		if err == nil {
			fmt.Println(v)
		}
	case "redo":
		err = goose.Redo(db, *migrationsDir)
	case "reset":
		err = goose.Reset(db, *migrationsDir)
	default:
		log.Fatalf("unsupported command: %s", *command)
	}

	if err != nil {
		log.Fatalf("goose %s failed: %v", *command, err)
	}

	fmt.Fprintf(os.Stdout, "goose %s completed successfully\n", *command)
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
