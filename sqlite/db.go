package sqlite

import "database/sql"

type DB struct {
	DB     *sql.DB
	Schema string
}
