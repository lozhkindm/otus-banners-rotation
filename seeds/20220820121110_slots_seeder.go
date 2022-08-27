package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3" //nolint:typecheck
)

func init() {
	goose.AddMigration(upSlotsSeeder, downSlotsSeeder) //nolint:typecheck
}

func upSlotsSeeder(tx *sql.Tx) error {
	query := `
insert into slots (name, created_at) VALUES
('header', NOW()),
('sidebar', NOW()),
('footer', NOW())
;`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}

func downSlotsSeeder(_ *sql.Tx) error {
	return nil
}
