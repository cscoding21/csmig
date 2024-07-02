package persistence

import (
	"fmt"

	"github.com/cscoding21/csmig/shared"
	"github.com/surrealdb/surrealdb.go"
)

var SurrealDBStrategy = DatabaseStrategy{
	Name: "surrealdb",
	DBConfig: DatabaseConfig{
		Name:      "surrealdb",
		Host:      "localhost",
		Port:      9999,
		User:      "root",
		Password:  "root",
		Database:  "test",
		Namespace: "test",
	},
	EnsureInfrastructure: func(config DatabaseConfig) error {
		db, err := getSurrealDB(config)
		if err != nil {
			panic(err)
		}

		defineSQL := `
		DEFINE TABLE IF NOT EXISTS csmig_versions SCHEMAFULL;
		DEFINE FIELD IF NOT EXISTS name ON TABLE csmig_versions TYPE string;
		DEFINE FIELD IF NOT EXISTS applied_on ON TABLE csmig_versions TYPE datetime DEFAULT time::now();
		DEFINE INDEX csmig_versions_name_unique ON TABLE csmig_versions COLUMNS name UNIQUE;
		`

		_, err = db.Query(defineSQL, nil)
		if err != nil {
			panic(err)
		}

		return nil
	},
	ApplyMigration: func(config DatabaseConfig, name string) error {
		db, err := getSurrealDB(config)
		if err != nil {
			panic(err)
		}

		applySQL := fmt.Sprintf(`INSERT INTO %s (name) VALUES ($name);`, VersionTableName)

		_, err = db.Query(applySQL, map[string]interface{}{
			"name": name,
		})
		if err != nil {
			panic(err)
		}

		return nil
	},
	FindAppliedMigrations: func(config DatabaseConfig) ([]shared.AppliedMigration, error) {
		db, err := getSurrealDB(config)
		if err != nil {
			panic(err)
		}

		applySQL := fmt.Sprintf(`SELECT * FROM %s ORDER BY applied_on ASC;`, VersionTableName)
		migrationData, err := db.Query(applySQL, nil)
		if err != nil {
			panic(err)
		}

		appliedMigraitons, err := surrealdb.SmartUnmarshal[[]shared.AppliedMigration](migrationData, err)
		if err != nil {
			panic(err)
		}

		return appliedMigraitons, nil
	},
	RollbackMigration: func(config DatabaseConfig, name string) error {
		db, err := getSurrealDB(config)
		if err != nil {
			panic(err)
		}

		applySQL := fmt.Sprintf(`DELETE FROM %s where name = $name;`, VersionTableName)

		_, err = db.Query(applySQL, map[string]interface{}{
			"name": name,
		})
		if err != nil {
			panic(err)
		}

		return nil
	},
	ResetMigrations: func(config DatabaseConfig) error {
		db, err := getSurrealDB(config)
		if err != nil {
			panic(err)
		}

		applySQL := fmt.Sprintf(`DELETE FROM %s;`, VersionTableName)

		_, err = db.Query(applySQL, nil)
		if err != nil {
			panic(err)
		}

		return nil
	},
}

func getSurrealDB(config DatabaseConfig) (*surrealdb.DB, error) {
	db, err := surrealdb.New(fmt.Sprintf("ws://%s:%v/rpc", config.Host, config.Port))
	if err != nil {
		panic(err)
	}

	// Sign in
	if _, err = db.Signin(map[string]string{
		"user": config.User,
		"pass": config.Password,
	}); err != nil {
		return nil, err
	}

	// Select namespace and database
	if _, err = db.Use(config.Namespace, config.Database); err != nil {
		return nil, err
	}

	return db, nil
}
