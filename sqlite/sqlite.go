// Package sqlite has helpers for working with sqlite.
package sqlite

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/url"

	"github.com/harrybrwn/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

// Open will open a database with an explicitly passed config.
func Open(location string, config *Config) (*sql.DB, error) {
	if config == nil {
		panic("sqlite: *Config is required to open a database")
	}
	config.logger = slog.New(slog.DiscardHandler)
	query, err := config.query()
	if err != nil {
		return nil, err
	}
	uri := url.URL{
		Scheme:   "file",
		Opaque:   location,
		RawQuery: query.Encode(),
	}
	return open(&uri, config)
}

// OpenURI will accept a URL and options to open a new database.
func OpenURI(uri *url.URL, opts ...Option) (*sql.DB, error) {
	var config Config
	config.logger = slog.New(slog.DiscardHandler)
	for _, o := range opts {
		o(&config)
	}
	query, err := config.query()
	if err != nil {
		return nil, err
	}
	uri.RawQuery = query.Encode()
	return open(uri, &config)
}

// File will open a database file.
func File(location string, opts ...Option) (*sql.DB, error) {
	var config Config
	config.logger = slog.New(slog.DiscardHandler)
	for _, o := range opts {
		o(&config)
	}
	query, err := config.query()
	if err != nil {
		return nil, err
	}
	uri := url.URL{
		Scheme:   "file",
		Opaque:   location,
		RawQuery: query.Encode(),
	}
	return open(&uri, &config)
}

// InMemory will open an in-memory database.
func InMemory(opts ...Option) (*sql.DB, error) {
	var config Config
	config.logger = slog.New(slog.DiscardHandler)
	for _, o := range opts {
		o(&config)
	}
	query, err := config.query()
	if err != nil {
		return nil, err
	}
	uri := url.URL{
		Opaque:   ":memory:",
		RawQuery: query.Encode(),
	}
	return open(&uri, &config)
}

func open(uri *url.URL, config *Config) (*sql.DB, error) {
	source := uri.String()
	if config != nil && config.Debug {
		config.logger.Debug("sql.Open", "driver", "sqlite3", "source", source)
	}
	db, err := sql.Open("sqlite3", source)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = config.pragmas(db); err != nil {
		return nil, err
	}
	return db, nil
}

func ListTablesNames(db db.DB) ([]string, error) {
	rows, err := db.QueryContext(context.Background(), `SELECT tbl_name FROM sqlite_schema WHERE type = 'table'`)
	if err != nil {
		return nil, err
	}
	defer closeRows(db, rows)
	return scanStrings(rows)
}

// SchemaType is a string enum for the 'sqlite_schema' type column.
type SchemaType string

const (
	SchemaTypeTable   SchemaType = "table"
	SchemaTypeIndex   SchemaType = "index"
	SchemaTypeView    SchemaType = "view"
	SchemaTypeTrigger SchemaType = "trigger"
)

// Schema represents a row from the 'sqlite_schema' table.
type Schema struct {
	Type      SchemaType
	Name      string
	TableName string
	Rootpage  string
	SQL       sql.NullString
}

// GetSchemas queries the 'sqlite_schema' table.
func GetSchemas(db db.DB) ([]Schema, error) {
	rows, err := db.QueryContext(
		context.Background(),
		`SELECT type,name,tbl_name,rootpage,sql FROM sqlite_schema`,
	)
	if err != nil {
		return nil, err
	}
	defer closeRows(db, rows)
	return autoscanRows[Schema](rows)
}

func exec(c *Config, database *sql.DB, query string, args ...any) error {
	if c.Debug {
		c.logger.Debug("EXEC", "query", query, "args", args)
	}
	_, err := database.Exec(query, args...)
	return errors.WithStack(err)
}

func closeRows(database any, rows io.Closer) {
	err := rows.Close()
	if err != nil {
		var logger *slog.Logger
		if loggerdb, ok := database.(interface{ Logger() *slog.Logger }); ok {
			logger = loggerdb.Logger()
		}
		if logger == nil {
			logger = slog.Default()
		}
		logger.Error("failed to close database rows", "error", err)
	}
}
