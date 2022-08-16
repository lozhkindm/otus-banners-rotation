package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddBanner(t *testing.T) {
	defer truncateTables()
	slotID := 1
	bannerID := 1

	inRotation, err := isBannerInRotation(slotID, bannerID)
	require.NoError(t, err)
	require.False(t, inRotation)

	statusCode, err := addBannerToRotation(slotID, bannerID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	inRotation, err = isBannerInRotation(slotID, bannerID)
	require.NoError(t, err)
	require.True(t, inRotation)
}

func TestRemoveBanner(t *testing.T) {
	defer truncateTables()
	slotID := 1
	bannerID := 1

	statusCode, err := addBannerToRotation(slotID, bannerID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	inRotation, err := isBannerInRotation(slotID, bannerID)
	require.NoError(t, err)
	require.True(t, inRotation)

	statusCode, err = removeBannerFromRotation(slotID, bannerID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	inRotation, err = isBannerInRotation(slotID, bannerID)
	require.NoError(t, err)
	require.False(t, inRotation)
}

func TestClickBanner(t *testing.T) {
	defer truncateTables()
	slotID := 1
	bannerID := 1
	usergroupID := 1

	statusCode, err := addBannerToRotation(slotID, bannerID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	inRotation, err := isBannerInRotation(slotID, bannerID)
	require.NoError(t, err)
	require.True(t, inRotation)

	clicks, err := countBannerClicks(slotID, bannerID, usergroupID)
	require.NoError(t, err)
	require.Equal(t, 0, clicks)

	statusCode, err = clickOnBanner(slotID, bannerID, usergroupID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	clicks, err = countBannerClicks(slotID, bannerID, usergroupID)
	require.NoError(t, err)
	require.Equal(t, 1, clicks)
}

func TestPickBanner(t *testing.T) {
	defer truncateTables()
	slotID := 1
	bannerID := 1
	usergroupID := 1

	statusCode, err := addBannerToRotation(slotID, bannerID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	inRotation, err := isBannerInRotation(slotID, bannerID)
	require.NoError(t, err)
	require.True(t, inRotation)

	impressions, err := countBannerImpressions(slotID, bannerID, usergroupID)
	require.NoError(t, err)
	require.Equal(t, 0, impressions)

	statusCode, bannerID, err = pickBanner(slotID, usergroupID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, 1, bannerID)

	impressions, err = countBannerImpressions(slotID, bannerID, usergroupID)
	require.NoError(t, err)
	require.Equal(t, 1, impressions)
}

func countBannerClicks(slotID, bannerID, usergroupID int) (int, error) {
	args := map[string]interface{}{
		"slot_id":      bannerID,
		"banner_id":    slotID,
		"usergroup_id": usergroupID,
	}
	query := `
		select count(*) as clicks
		from clicks 
		where slot_id = :slot_id and banner_id = :banner_id and usergroup_id = :usergroup_id
	;`
	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var count int
	if err := stmt.QueryRow(args).Scan(&count); err != nil {
		return 0, err
	}
	return count, err
}

func countBannerImpressions(slotID, bannerID, usergroupID int) (int, error) {
	args := map[string]interface{}{
		"slot_id":      bannerID,
		"banner_id":    slotID,
		"usergroup_id": usergroupID,
	}
	query := `
		select count(*) as impressions
		from impressions 
		where slot_id = :slot_id and banner_id = :banner_id and usergroup_id = :usergroup_id
	;`
	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var count int
	if err := stmt.QueryRow(args).Scan(&count); err != nil {
		return 0, err
	}
	return count, err
}

func isBannerInRotation(slotID, bannerID int) (bool, error) {
	args := map[string]interface{}{
		"slot_id":   bannerID,
		"banner_id": slotID,
	}
	query := "select exists(select 1 from rotations where slot_id = :slot_id and banner_id = :banner_id)"
	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var exists bool
	if err := stmt.QueryRow(args).Scan(&exists); err != nil {
		return false, err
	}
	return exists, err
}

func addBannerToRotation(slotID, bannerID int) (int, error) {
	req := addBannerRequest{
		BannerID: bannerID,
		SlotID:   slotID,
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(req); err != nil {
		return 0, err
	}

	res, err := http.Post(host+addBannerUrl, "application/json", &body)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return res.StatusCode, nil
}

func removeBannerFromRotation(slotID, bannerID int) (int, error) {
	req := removeBannerRequest{
		BannerID: bannerID,
		SlotID:   slotID,
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(req); err != nil {
		return 0, err
	}

	request, err := http.NewRequest(http.MethodDelete, host+removeBannerUrl, &body)
	if err != nil {
		return 0, err
	}

	request.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return res.StatusCode, nil
}

func clickOnBanner(slotID, bannerID, usergroupID int) (int, error) {
	req := clickBannerRequest{
		BannerID:    bannerID,
		SlotID:      slotID,
		UsergroupID: usergroupID,
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(req); err != nil {
		return 0, err
	}

	res, err := http.Post(host+clickBannerUrl, "application/json", &body)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return res.StatusCode, nil
}

func pickBanner(slotID, usergroupID int) (int, int, error) {
	params := url.Values{}
	params.Add("slotId", strconv.Itoa(slotID))
	params.Add("usergroupId", strconv.Itoa(usergroupID))

	res, err := http.Get(host + pickBannerUrl + "?" + params.Encode())
	if err != nil {
		fmt.Println(123123132)
		return 0, 0, err
	}
	defer res.Body.Close()

	var response pickBannerResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return 0, 0, err
	}
	return res.StatusCode, response.BannerID, nil
}

func truncateTables() {
	if _, err := db.Exec("truncate table rotations;"); err != nil {
		panic(err)
	}
	if _, err := db.Exec("truncate table impressions;"); err != nil {
		panic(err)
	}
	if _, err := db.Exec("truncate table clicks;"); err != nil {
		panic(err)
	}
}
