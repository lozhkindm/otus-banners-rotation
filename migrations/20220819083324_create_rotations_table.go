package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3" //nolint:typecheck
)

func init() {
	goose.AddMigration(upCreateRotationsTable, downCreateRotationsTable) //nolint:typecheck
}

func upCreateRotationsTable(tx *sql.Tx) error {
	query := `
create table if not exists rotations
(
    slot_id    int       not null
        constraint rotations_slots_id_fk
            references slots
            on update cascade on delete cascade,
    banner_id  int       not null
        constraint rotations_banners_id_fk
            references banners
            on update cascade on delete cascade,
    created_at timestamp not null,
    constraint rotations_pk
        primary key (slot_id, banner_id)
);
`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}

func downCreateRotationsTable(tx *sql.Tx) error {
	query := "drop table if exists rotations;"
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}
