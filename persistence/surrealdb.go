package persistence

import (
	"fmt"

	"github.com/cscoding21/csmig/shared"
	"github.com/surrealdb/surrealdb.go"
)

var _conn *surrealdb.DB

var SurrealDBStrategy = shared.DatabaseStrategy{
	Name: "surrealdb",
	DBConfig: shared.DatabaseConfig{
		Name:      "surrealdb",
		Host:      "localhost",
		Port:      9999,
		User:      "root",
		Password:  "root",
		Database:  "test",
		Namespace: "test",
	},
	EnsureInfrastructure: func(config shared.DatabaseConfig) error {
		db, err := GetConnection(config)
		if err != nil {
			panic(err)
		}

		defineSQL := fmt.Sprintf(`
		DEFINE TABLE IF NOT EXISTS %s SCHEMAFULL;
		DEFINE FIELD IF NOT EXISTS name ON TABLE %s TYPE string;
		DEFINE FIELD IF NOT EXISTS description ON TABLE %s TYPE string;
		DEFINE FIELD IF NOT EXISTS applied_on ON TABLE %s TYPE datetime DEFAULT time::now();
		DEFINE INDEX %s_name_unique ON TABLE %s COLUMNS name UNIQUE;
		`, VersionTableName, VersionTableName, VersionTableName, VersionTableName, VersionTableName, VersionTableName)
		_, err = db.Query(defineSQL, nil)
		if err != nil {
			return err
		}

		return nil
	},
	ApplyMigration: func(config shared.DatabaseConfig, name string, description string) error {
		db, err := GetConnection(config)
		if err != nil {
			return err
		}

		applySQL := fmt.Sprintf(`INSERT INTO %s (name, description) VALUES ($name, $description);`, VersionTableName)

		_, err = db.Query(applySQL, map[string]interface{}{
			"name":        name,
			"description": description,
		})
		if err != nil {
			return err
		}

		return nil
	},
	FindAppliedMigrations: func(config shared.DatabaseConfig) ([]shared.AppliedMigration, error) {
		db, err := GetConnection(config)
		if err != nil {
			return nil, err
		}

		applySQL := fmt.Sprintf(`SELECT * FROM %s ORDER BY applied_on ASC;`, VersionTableName)
		migrationData, err := db.Query(applySQL, nil)
		if err != nil {
			return nil, err
		}

		appliedMigraitons, err := surrealdb.SmartUnmarshal[[]shared.AppliedMigration](migrationData, err)
		if err != nil {
			return nil, err
		}

		return appliedMigraitons, nil
	},
	RollbackMigration: func(config shared.DatabaseConfig, name string) error {
		db, err := GetConnection(config)
		if err != nil {
			return err
		}

		applySQL := fmt.Sprintf(`DELETE FROM %s where name = $name;`, VersionTableName)

		_, err = db.Query(applySQL, map[string]interface{}{
			"name": name,
		})
		if err != nil {
			return err
		}

		return nil
	},
	ResetMigrations: func(config shared.DatabaseConfig) error {
		db, err := GetConnection(config)
		if err != nil {
			return err
		}

		applySQL := fmt.Sprintf(`DELETE FROM %s;`, VersionTableName)

		_, err = db.Query(applySQL, nil)
		if err != nil {
			return err
		}

		return nil
	},
	Exec: func(config shared.DatabaseConfig, sql string, params map[string]interface{}) error {
		db, err := GetConnection(config)
		if err != nil {
			return err
		}

		_, err = db.Query(sql, params)
		if err != nil {
			return err
		}

		return nil
	},
}

func GetConnection(config shared.DatabaseConfig) (*surrealdb.DB, error) {
	if _conn != nil {
		return _conn, nil
	}

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

	_conn = db

	return _conn, nil
}
