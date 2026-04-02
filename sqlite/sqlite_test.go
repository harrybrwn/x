package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/harrybrwn/db"
	"github.com/matryer/is"
	"github.com/mattn/go-sqlite3"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

func TestConfig(t *testing.T) {
	type testcase struct {
		name     string
		u        url.URL
		c        *Config
		query    func(*is.I, url.Values)
		queryErr error
	}

	for _, tc := range []testcase{
		{
			"cache_shared",
			url.URL{
				Scheme:   "file",
				Path:     filepath.Join(t.TempDir(), "test.sqlite"),
				RawQuery: must((&Config{Cache: CacheModePrivate}).query()).Encode(),
			},
			nil,
			func(*is.I, url.Values) {},
			nil,
		},

		{
			"readonly",
			url.URL{Opaque: ":memory:"},
			&Config{ReadOnly: true, Cache: CacheModeShared},
			func(is *is.I, q url.Values) {
				is.Equal(q.Get("mode"), "ro")
				is.Equal(q.Get("immutable"), "true")
				is.Equal(q.Get("cache"), "shared")
			},
			nil,
		},

		{
			"invalid-cache-mode",
			url.URL{Opaque: ":memory:"},
			&Config{Cache: CacheMode(99)},
			func(i *is.I, v url.Values) {},
			ErrInvalidCacheMode,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			var (
				d   *sql.DB
				q   url.Values
				err error
			)

			if tc.c == nil {
				d, err = OpenURI(&tc.u)
			} else {
				q, err = tc.c.query()
				if tc.queryErr != nil {
					is.True(errors.Is(err, tc.queryErr))
					return
				}
				is.NoErr(err)
				tc.query(is, q)
				d, err = Open(tc.u.String(), tc.c)
			}
			is.NoErr(err)
			defer closeDB(d)
		})
	}
}

func TestWierdConfig(t *testing.T) {
	is := is.New(t)
	var c *Config
	q, err := c.query()
	is.NoErr(err)
	is.Equal(q, url.Values{})
}

func TestOpen(t *testing.T) {
	is := is.New(t)
	d, err := Open(
		filepath.Join(t.TempDir(), "test.sqlite"),
		&Config{
			JournalMode: "TRUNCATE",
		},
	)
	is.NoErr(err)
	defer closeDB(d)
	mode, err := GetJournalMode(db.Simple(d))
	is.NoErr(err)
	is.Equal(mode, "truncate")
}

func TestInMemory(t *testing.T) {
	is := is.New(t)
	d, err := InMemory(ReadOnly)
	is.NoErr(err)
	defer closeDB(d)
	mode, err := GetJournalMode(db.Simple(d))
	is.NoErr(err)
	is.Equal(mode, "memory")
}

func TestOpenURI(t *testing.T) {
	is := is.New(t)
	d, err := OpenURI(&url.URL{Opaque: ":memory:"}, WithPragma(PragmaForeignKeys, 1), ReadOnly)
	is.NoErr(err)
	defer closeDB(d)
	on, err := GetPragma[bool](db.Simple(d), PragmaForeignKeys)
	is.NoErr(err)
	is.True(on)
}

func TestPragmas(t *testing.T) {
	is := is.New(t)
	uri := url.URL{
		Scheme: "file",
		Opaque: filepath.Join(t.TempDir(), "test.db"),
	}
	d, err := open(&uri, new(Config))
	is.NoErr(err)
	defer closeDB(d)
	db := db.Simple(d)
	mode, err := GetJournalMode(db)
	is.NoErr(err)
	is.Equal(mode, "delete") // delete should be the default
	c := Config{
		JournalMode:   "WAL",
		WalCheckpoint: ptr(7), // TODO test the results of this
	}
	WithPragma(PragmaSynchronous, SynchronousExtra)(&c)
	WithPragma(PragmaCacheSize, 69)(&c)
	is.Equal(c.Pragmas, map[string]any{
		"cache_size":      69,
		PragmaSynchronous: SynchronousExtra,
	})
	err = c.pragmas(d)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	is.NoErr(err)
	mode, err = GetJournalMode(db)
	is.NoErr(err)
	is.Equal(mode, "wal")
	sync, err := GetPragmaSynchronous(db)
	is.NoErr(err)
	is.Equal(sync, SynchronousExtra)
	cacheSize, err := GetPragmaCacheSize(db)
	is.NoErr(err)
	is.Equal(cacheSize, int64(69))
	databases, err := GetPragmaDatabaseList(db)
	is.NoErr(err)
	is.Equal(len(databases), 1)
	is.Equal(databases[0].Index, 0)
	is.Equal(databases[0].Name, "main")
	is.Equal(databases[0].Location, uri.Opaque)
	_, _, _, err = GetWalCheckpoint(db) // TODO check values
	is.NoErr(err)
}

func TestListTables(t *testing.T) {
	is := is.New(t)
	d, err := File(
		filepath.Join(t.TempDir(), "test.sqlite"),
		Cache(CacheModePrivate),
		WithPragma(PragmaForeignKeys, 0),
		JournalMode("WAL"),
		WalCheckpoint(0),
		debug,
	)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	is.NoErr(err)
	defer closeDB(d)
	on, err := GetPragma[bool](db.Simple(d), PragmaForeignKeys)
	is.NoErr(err)
	is.True(!on)
	_, err = d.Exec(`CREATE TABLE testing_table (
		name varchar,
		number INT
	)`)
	is.NoErr(err)
	names, err := ListTablesNames(db.Simple(d))
	is.NoErr(err)
	is.Equal(names, []string{"testing_table"})
}

func TestGetSchemas(t *testing.T) {
	is := is.New(t)
	d, err := InMemory()
	is.NoErr(err)
	defer closeDB(d)
	table1SQL := `CREATE TABLE test_table (
		identifier INTEGER PRIMARY KEY,
		name blob,
		count INTEGER,
		time TIMESTAMP NOT NULL DEFAULT 99
	)`
	table2SQL := `CREATE TABLE test_table2 (
		name TEXT PRIMARY KEY,
		t1_id INTEGER NOT NULL,
	    FOREIGN KEY (t1_id) REFERENCES "test_table" (identifier) ON DELETE CASCADE
	)`
	_, err = d.Exec(table1SQL)
	is.NoErr(err)
	_, err = d.Exec(table2SQL)
	is.NoErr(err)
	schemas, err := GetSchemas(db.Simple(d))
	is.NoErr(err)
	is.Equal(len(schemas), 3)
	is.Equal(schemas[0], Schema{Type: SchemaTypeTable, Name: "test_table", TableName: "test_table", Rootpage: "2", SQL: sql.NullString{String: table1SQL, Valid: true}})
	is.Equal(schemas[1], Schema{Type: SchemaTypeTable, Name: "test_table2", TableName: "test_table2", Rootpage: "3", SQL: sql.NullString{String: table2SQL, Valid: true}})
	is.Equal(schemas[2], Schema{Type: SchemaTypeIndex, Name: "sqlite_autoindex_test_table2_1", TableName: "test_table2", Rootpage: "4", SQL: sql.NullString{Valid: false}})
}

func TestTableInfo(t *testing.T) {
	is := is.New(t)
	d, err := InMemory()
	is.NoErr(err)
	defer closeDB(d)
	_, err = d.Exec(`CREATE TABLE test_table (
		identifier INTEGER PRIMARY KEY,
		name blob,
		count INTEGER,
		time TIMESTAMP NOT NULL DEFAULT 99
	)`)
	is.NoErr(err)
	columns, err := GetTableInfo(db.Simple(d), "test_table")
	is.NoErr(err)
	is.Equal(len(columns), 4)
	is.Equal(columns[0], TableInfo{
		CID: 0, Name: "identifier", Type: "INTEGER", NotNull: false,
		Default: sql.NullString{String: "", Valid: false}, PrimaryKey: true,
	})
	is.Equal(columns[1], TableInfo{
		CID: 1, Name: "name", Type: "BLOB", NotNull: false,
		Default: sql.NullString{String: "", Valid: false}, PrimaryKey: false,
	})
	is.Equal(columns[2], TableInfo{
		CID: 2, Name: "count", Type: "INTEGER", NotNull: false,
		Default: sql.NullString{String: "", Valid: false}, PrimaryKey: false,
	})
	is.Equal(columns[3], TableInfo{
		CID: 3, Name: "time", Type: "TIMESTAMP", NotNull: true,
		Default: sql.NullString{String: "99", Valid: true}, PrimaryKey: false,
	})
}

func TestQueryOnly(t *testing.T) {
	is := is.New(t)
	d, err := File(
		filepath.Join(t.TempDir(), "test.sqlite"),
		WithPragma(PragmaQueryOnly, 1),
		Logger(slog.New(slog.DiscardHandler)),
		WithSynchronous(SynchronousExtra),
		Debug(true),
	)
	is.NoErr(err)
	defer closeDB(d)
	_, err = d.Exec(`CREATE TABLE t(a TEXT, b INTEGER)`)
	is.True(err != nil)
	e, ok := err.(sqlite3.Error)
	is.True(ok)
	is.Equal(e.Code, sqlite3.ErrReadonly)
	is.Equal(e.Error(), "attempt to write a readonly database")
}

func debug(c *Config) { c.Debug = true }

func must[T any](v T, e error) T {
	if e != nil {
		panic(e)
	}
	return v
}

func ptr[T any](v T) *T { return &v }

func closeDB(database io.Closer) {
	err := database.Close()
	if err != nil {
		logger := slog.Default()
		if loggerdb, ok := database.(interface{ Logger() *slog.Logger }); ok {
			logger = loggerdb.Logger()
		}
		logger.Error("failed to close database", "error", err)
	}
}
