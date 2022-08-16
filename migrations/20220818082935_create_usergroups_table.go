package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3" //nolint:typecheck
)

func init() {
	goose.AddMigration(upCreateUsergroupsTable, downCreateUsergroupsTable) //nolint:typecheck
}

func upCreateUsergroupsTable(tx *sql.Tx) error {
	query := `
create table if not exists usergroups
(
    id         serial constraint usergroups_pk primary key,
    name       varchar   not null,
    created_at timestamp not null
);
`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}

func downCreateUsergroupsTable(tx *sql.Tx) error {
	query := `drop table if exists usergroups;`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}
