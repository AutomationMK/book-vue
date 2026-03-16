package driver

import (
	"context"

	_ "github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB holds the database connection pool
type DB struct {
	SQL *pgx.Conn
}

var dbConn = &DB{}

// ConnectSQL creates dabase pool for postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	dbConn.SQL = d

	return dbConn, nil
}

/* testDb tries to ping the database
func testDb(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}
*/

// NewDatabase creates a new database for the application
func NewDatabase(dsn string) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
