package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	//Драйвер для миграций на sqlite3
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	//Драйвер для получения миграций из файлов
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		storagePath     string
		migrationsPath  string
		migrationsTable string
	)

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations")
	flag.Parse()

	if storagePath == "" || migrationsPath == "" {
		panic("storage-path and migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")
			return
		}
		panic(err)
	}
	fmt.Println("migrations applied successfully")
}
