package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3" //nolint:typecheck
)

func init() {
	goose.AddMigration(upCreateSlotsTable, downCreateSlotsTable) //nolint:typecheck
}

func upCreateSlotsTable(tx *sql.Tx) error {
	query := `
create table if not exists slots
(
    id         serial constraint slots_pk primary key,
    name       varchar   not null,
    created_at timestamp not null
);
`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}

func downCreateSlotsTable(tx *sql.Tx) error {
	query := "drop table if exists slots;"
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}
