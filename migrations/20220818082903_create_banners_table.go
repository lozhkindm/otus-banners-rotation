package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3" //nolint:typecheck
)

func init() {
	goose.AddMigration(upCreateBannersTable, downCreateBannersTable) //nolint:typecheck
}

func upCreateBannersTable(tx *sql.Tx) error {
	query := `
create table if not exists banners
(
    id         serial constraint banners_pk primary key,
    name       varchar   not null,
    created_at timestamp not null
);
`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}

func downCreateBannersTable(tx *sql.Tx) error {
	query := "drop table if exists banners;"
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}
