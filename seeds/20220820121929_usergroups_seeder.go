package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3" //nolint:typecheck
)

func init() {
	goose.AddMigration(upUsergroupsSeeder, downUsergroupsSeeder) //nolint:typecheck
}

func upUsergroupsSeeder(tx *sql.Tx) error {
	query := `
insert into usergroups (name, created_at) VALUES
('children', NOW()),
('teenagers', NOW()),
('old people', NOW())
;`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}

func downUsergroupsSeeder(_ *sql.Tx) error {
	return nil
}
