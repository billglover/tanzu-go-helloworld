package visits

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/url"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var efs embed.FS

const version = 1

type VisitorLog struct {
	db *sql.DB
}

func NewLog(pgURL *url.URL) (*VisitorLog, error) {
	db, err := sql.Open("pgx", pgURL.String())
	if err != nil {
		return nil, fmt.Errorf("unable to open DB: %w", err)
	}

	err = validateSchema(db)
	if err != nil {
		return nil, err
	}

	return &VisitorLog{
		db: db,
	}, nil
}

func (vl *VisitorLog) NewVisit(ctx context.Context) (string, error) {
	row := vl.db.QueryRowContext(ctx, "INSERT INTO visits DEFAULT VALUES RETURNING id;")
	var id string
	err := row.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("unable to scan row: %w", err)
	}
	return id, nil
}

// Migrate migrates the Postgres schema to the current version.
func validateSchema(db *sql.DB) error {

	d, err := iofs.New(efs, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	targetInstance, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", targetInstance)
	if err != nil {
		return err
	}

	err = m.Migrate(version) // current version
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return d.Close()
}
