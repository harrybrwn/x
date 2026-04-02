package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/harrybrwn/db"
	"github.com/pkg/errors"
)

const (
	// Settings pragmas

	PragmaSynchronous            = "synchronous"
	PragmaJournalMode            = "journal_mode"
	PragmaWalCheckpoint          = "wal_checkpoint"
	PragmaCacheSize              = "cache_size"
	PragmaApplicationID          = "application_id"
	PragmaAutoVacuum             = "auto_vacuum"
	PragmaAutomaticIndex         = "automatic_index"
	PragmaDataVersion            = "data_version"
	PragmaForeignKeys            = "foreign_keys"
	PragmaEncoding               = "encoding"
	PragmaIgnoreCheckConstraints = "ignore_check_constraints"
	PragmaLockingMode            = "locking_mode"
	PragmaQueryOnly              = "query_only"

	// Info pragmas

	PragmaDatabaseList   = "database_list"
	PragmaIndexList      = "index_list"
	PragmaForeignKeyList = "foreign_key_list"
	PragmaFunctionList   = "function_list"
	PragmaTableInfo      = "table_info"
	PragmaIndexInfo      = "index_info"
	PragmaPragmaList     = "pragma_list"
	PragmaModuleList     = "module_list"
)

// Synchronous is an enum for sqlite synchronous pragma settings. See
// https://www.sqlite.org/pragma.html
type Synchronous uint8

const (
	SynchronousOff Synchronous = iota
	SynchronousNormal
	SynchronousFull
	SynchronousExtra
)

func (s Synchronous) String() string {
	switch s {
	case SynchronousOff:
		return "OFF"
	case SynchronousNormal:
		return "NORMAL"
	case SynchronousFull:
		return "FULL"
	case SynchronousExtra:
		return "EXTRA"
	}
	return ""
}

func GetPragma[T any](database db.DB, name string) (T, error) {
	var v T
	return v, getPragma(database, name, &v)
}

func GetJournalMode(database db.DB) (string, error) {
	return GetPragma[string](database, PragmaJournalMode)
}

func GetPragmaSynchronous(database db.DB) (Synchronous, error) {
	return GetPragma[Synchronous](database, PragmaSynchronous)
}

func GetPragmaCacheSize(database db.DB) (int64, error) {
	return GetPragma[int64](database, PragmaCacheSize)
}

func GetWalCheckpoint(database db.DB) (int, int, int, error) {
	var a, b, c int
	return a, b, c, getPragma(database, PragmaWalCheckpoint, &a, &b, &c)
}

type DatabaseList struct {
	Index    int
	Name     string
	Location string
}

func GetPragmaDatabaseList(database db.DB) ([]DatabaseList, error) {
	rows, err := database.QueryContext(
		context.Background(),
		`PRAGMA `+PragmaDatabaseList,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer closeRows(database, rows)
	return autoscanRows[DatabaseList](rows)
}

// TableInfo holds a row from the 'table_info' pragma.
type TableInfo struct {
	CID        int
	Name       string
	Type       string
	NotNull    bool
	Default    sql.NullString
	PrimaryKey bool
}

// GetTableInfo gets info for each column using the 'table_info' pragma.
func GetTableInfo(database db.DB, table string) ([]TableInfo, error) {
	rows, err := database.QueryContext(
		context.Background(),
		fmt.Sprintf("PRAGMA %s(%s)", PragmaTableInfo, table))
	if err != nil {
		return nil, err
	}
	defer closeRows(database, rows)
	return autoscanRows[TableInfo](rows)
}

type ForeignKey struct {
	ID       int
	Seq      int
	Table    string
	From     string
	To       string
	OnUpdate string
	OnDelete string
	Match    string
}

func ForeignKeyList(database db.DB, table string) ([]ForeignKey, error) {
	rows, err := database.QueryContext(context.Background(), fmt.Sprintf("PRAGMA %s(%s)", PragmaForeignKeyList, table))
	if err != nil {
		return nil, err
	}
	defer closeRows(database, rows)
	return autoscanRows[ForeignKey](rows)
}

type IndexInfo struct {
	SeqNumber int
	CID       int
	Name      string
}

func GetIndexInfo(database db.DB, name string) (IndexInfo, error) {
	var i IndexInfo
	return i, getPragma(database, fmt.Sprintf("%s(%s)", PragmaIndexInfo, name), &i.SeqNumber, &i.CID, &i.Name)
}

type IndexXInfo struct {
	SeqNumber int
	CID       int
	Name      string
	Desc      bool
	Coll      string
	Key       int
}

func ModuleList(database db.DB) ([]string, error) {
	return getPragmaStrings(database, PragmaModuleList)
}

func getPragmaStrings(database db.DB, pragma string) ([]string, error) {
	rows, err := database.QueryContext(context.Background(), `PRAGMA `+pragma)
	if err != nil {
		return nil, err
	}
	defer closeRows(database, rows)
	return scanStrings(rows)
}

func getPragma(database db.DB, name string, dst ...any) error {
	rows, err := database.QueryContext(
		context.Background(),
		`PRAGMA `+name,
	)
	if err != nil {
		return err
	}
	return db.ScanOne(rows, dst...)
}

func pragma(c *Config, database *sql.DB, name string, value any) error {
	query := fmt.Sprintf("PRAGMA %s=%v", name, value)
	c.loggerOrDefault().Debug("executing pragma",
		"query", query)
	return exec(c, database, query)
}
