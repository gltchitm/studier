package db

import (
	"database/sql"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type parsedMigration struct {
	migrationId  int64
	migrationSql string
}

const migrationsPath string = "/studier/server/migrations"

func parseMigrations(lastMigrationId int64, migrationFilenames []string) ([]parsedMigration, error) {
	parsedMigrations := []parsedMigration{}

	for _, migrationFilename := range migrationFilenames {
		migrationId, err := strconv.ParseInt(strings.Split(migrationFilename, "-")[0], 10, 64)
		if err != nil {
			return nil, err
		}

		if migrationId > lastMigrationId {
			migrationFile, err := os.Open(migrationsPath + "/" + migrationFilename)
			if err != nil {
				return nil, err
			}

			migrationSql, err := io.ReadAll(migrationFile)
			if err != nil {
				return nil, err
			}

			parsedMigrations = append(parsedMigrations, parsedMigration{
				migrationId:  migrationId,
				migrationSql: string(migrationSql),
			})
		}
	}

	sort.Slice(parsedMigrations, func(i, j int) bool {
		return parsedMigrations[j].migrationId > parsedMigrations[i].migrationId
	})

	return parsedMigrations, nil
}

func migrateDb(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (migration_id INTEGER PRIMARY KEY);`)
	if err != nil {
		return err
	}

	var lastMigrationId int64
	err = db.QueryRow(`SELECT COALESCE(MAX(migration_id), -1) FROM migrations;`).Scan(&lastMigrationId)
	if err != nil {
		return err
	}

	dir, err := os.Open(migrationsPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	migrationFilenames, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}

	parsedMigrations, err := parseMigrations(lastMigrationId, migrationFilenames)
	if err != nil {
		return err
	}

	for _, parsedMigration := range parsedMigrations {
		_, err = db.Exec(parsedMigration.migrationSql)
		if err != nil {
			return err
		}

		_, err = db.Exec(`INSERT INTO migrations (migration_id) VALUES ($1)`, parsedMigration.migrationId)
		if err != nil {
			panic(err)
		}
	}

	return nil
}
