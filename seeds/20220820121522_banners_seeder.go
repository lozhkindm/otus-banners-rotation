package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3" //nolint:typecheck
)

func init() {
	goose.AddMigration(upBannersSeeder, downBannersSeeder) //nolint:typecheck
}

func upBannersSeeder(tx *sql.Tx) error {
	query := `
insert into banners (name, created_at) VALUES
('food', NOW()),
('toys', NOW()),
('medicines', NOW())
;`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}

func downBannersSeeder(_ *sql.Tx) error {
	return nil
}
