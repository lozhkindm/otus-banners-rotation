package sqlstorage

import (
	"context"
	"errors"
	"time"

	"github.com/lozhkindm/otus-banners-rotation/internal/multiarmedbandit"
	"github.com/lozhkindm/otus-banners-rotation/internal/storage"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const errCodeAlreadyExists = "23505"

var (
	errTimeoutClosingConn    = errors.New("timeout closing connection")
	errNoBannersForGivenSlot = errors.New("no banners for a given slot")
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Storage {
	return Storage{
		db: db,
	}
}

func (s Storage) AddBanner(ctx context.Context, bannerID, slotID int) error {
	args := map[string]interface{}{
		"banner_id": bannerID,
		"slot_id":   slotID,
	}
	query := `
		insert into rotations
		(slot_id, banner_id, created_at) values
		(:slot_id, :banner_id, NOW())
	;`

	_, err := s.db.NamedExecContext(ctx, query, args)

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case errCodeAlreadyExists:
			return nil
		default:
			return err
		}
	}

	return nil
}

func (s Storage) RemoveBanner(ctx context.Context, bannerID, slotID int) error {
	args := map[string]interface{}{
		"banner_id": bannerID,
		"slot_id":   slotID,
	}
	query := `delete from rotations where slot_id = :slot_id and banner_id = :banner_id;`

	if _, err := s.db.NamedExecContext(ctx, query, args); err != nil {
		return err
	}
	return nil
}

func (s Storage) ClickBanner(ctx context.Context, bannerID, slotID, usergroupID int) error {
	args := map[string]interface{}{
		"banner_id":    bannerID,
		"slot_id":      slotID,
		"usergroup_id": usergroupID,
	}
	query := `
		insert into clicks
		(slot_id, banner_id, usergroup_id, created_at) values
		(:slot_id, :banner_id, :usergroup_id, NOW())
	;`

	if _, err := s.db.NamedExecContext(ctx, query, args); err != nil {
		return err
	}
	return nil
}

func (s Storage) PickBanner(ctx context.Context, slotID, usergroupID int) (int, error) {
	args := map[string]interface{}{
		"slot_id":      slotID,
		"usergroup_id": usergroupID,
	}
	query := `
		select
			r.banner_id,
			(
				select count(*) from impressions i where i.banner_id = r.banner_id and i.usergroup_id = :usergroup_id
			) impressions,
			(
				select count(*) from clicks c where c.banner_id = r.banner_id and c.usergroup_id = :usergroup_id
			) clicks
		from rotations r
		where r.slot_id = :slot_id
	;`
	stmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, args)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	banners := make([]multiarmedbandit.Banner, 0)
	for rows.Next() {
		var bnr storage.BannerStatistics
		if err := rows.StructScan(&bnr); err != nil {
			return 0, err
		}
		banners = append(banners, &bnr)
	}

	if len(banners) == 0 {
		return 0, errNoBannersForGivenSlot
	}

	bannerID := multiarmedbandit.PickBanner(banners)
	if err := s.ImpressBanner(ctx, bannerID, slotID, usergroupID); err != nil {
		return 0, err
	}

	return bannerID, nil
}

func (s Storage) ImpressBanner(ctx context.Context, bannerID, slotID, usergroupID int) error {
	args := map[string]interface{}{
		"banner_id":    bannerID,
		"slot_id":      slotID,
		"usergroup_id": usergroupID,
	}
	query := `
		insert into impressions
		(slot_id, banner_id, usergroup_id, created_at) values
		(:slot_id, :banner_id, :usergroup_id, NOW())
	;`

	if _, err := s.db.NamedExecContext(ctx, query, args); err != nil {
		return err
	}
	return nil
}

func (s Storage) Connect(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s Storage) Close(ctx context.Context) error {
	var (
		ch             = make(chan error)
		newCtx, cancel = context.WithTimeout(ctx, 3*time.Second)
	)
	defer cancel()

	go func() {
		ch <- s.db.Close()
	}()

	select {
	case <-newCtx.Done():
		return errTimeoutClosingConn
	case <-ch:
		return nil
	}
}
