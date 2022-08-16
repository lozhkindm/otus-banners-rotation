package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3" //nolint:typecheck
)

func init() {
	goose.AddMigration(upCreateImpressionsTable, downCreateImpressionsTable) //nolint:typecheck
}

func upCreateImpressionsTable(tx *sql.Tx) error {
	query := `
create table if not exists impressions
(
    id           serial
        constraint impressions_pk
            primary key,
    slot_id      int       not null
        constraint impressions_slots_id_fk
            references slots
            on update cascade on delete cascade,
    banner_id    int       not null
        constraint impressions_banners_id_fk
            references banners
            on update cascade on delete cascade,
    usergroup_id int       not null
        constraint impressions_usergroups_id_fk
            references usergroups
            on update cascade on delete cascade,
    created_at   timestamp not null
);
`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}

func downCreateImpressionsTable(tx *sql.Tx) error {
	query := "drop table if exists impressions;"
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}
