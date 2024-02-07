package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/deividaspetraitis/go/database"
	"github.com/deividaspetraitis/go/errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// DB represents the database connection.
type DB struct {
	db     *sql.DB
	ctx    context.Context // background context
	cancel func()          // cancel background context

	// Datasource name
	DSN string

	// Path to database migrations
	MigrationSource string

	// Returns the current time. Defaults to time.Now()
	// Can be mocked for tests
	Now func() time.Time
}

// NewDB returns a new instance of DB associated with the given datasource name.
func NewDB(ctx context.Context, cfg *database.Config) *DB {
	db := &DB{
		DSN:             cfg.DSN(),
		MigrationSource: cfg.MigrationsSource,
		Now:             time.Now,
	}
	db.ctx, db.cancel = context.WithCancel(ctx)
	return db
}

// Open opens the database connection.
func (db *DB) Open() (err error) {
	// Ensure a DSN is set before attempting to open the database.
	if db.DSN == "" {
		return errors.New("dsn required")
	}

	// Connect to the database
	if db.db, err = sql.Open("postgres", db.DSN); err != nil {
		return errors.Wrap(err, "sql: open connection")
	}

	// Verify connection
	if err := db.db.Ping(); err != nil {
		return errors.Wrap(err, "sql: verify connection")
	}

	// Apply migrations
	if err := db.migrate(); err != nil {
		return errors.Wrap(err, "sql: migrate schema")
	}

	return nil
}

func (db *DB) migrate() error {
	driver, err := postgres.WithInstance(db.db, &postgres.Config{})
	if err != nil {
		return err
	}

	sourceURL := "file:///" + db.MigrationSource
	if len(db.MigrationSource) == 0 {
		sourceURL = "file://../../database/migrations"
	}

	instance, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return err
	}

	if err := instance.Up(); !errors.Equals(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	// Cancel background context.
	db.cancel()

	// Close database.
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Return wrapper Tx that includes the transaction start time.
	return &Tx{
		Tx:  tx,
		db:  db,
		now: db.Now().UTC().Truncate(time.Second),
	}, nil
}

// Tx wraps the SQL Tx object to provide a timestamp at the start of the transaction.
type Tx struct {
	*sql.Tx
	db  *DB
	now time.Time
}
