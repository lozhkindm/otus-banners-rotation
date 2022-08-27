package tests

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	addBannerUrl    = "/banner/rotation"
	removeBannerUrl = "/banner/rotation"
	clickBannerUrl  = "/banner/click"
	pickBannerUrl   = "/banner/pick"
)

var (
	db   *sqlx.DB
	host string
)

type addBannerRequest struct {
	BannerID int `json:"bannerId"`
	SlotID   int `json:"slotId"`
}

type removeBannerRequest struct {
	BannerID int `json:"bannerId"`
	SlotID   int `json:"slotId"`
}

type clickBannerRequest struct {
	BannerID    int `json:"bannerId"`
	SlotID      int `json:"slotId"`
	UsergroupID int `json:"usergroupId"`
}

type pickBannerResponse struct {
	BannerID int `json:"bannerId"`
}

func TestMain(m *testing.M) {
	fmt.Println(os.Getenv("POSTGRES_DSN"))
	var err error
	time.Local = nil
	db, err = sqlx.Open("pgx", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatal("failed to open pgx")
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %s", err)
	}
	host = os.Getenv("HTTP_HOST")
	os.Exit(m.Run())
}
