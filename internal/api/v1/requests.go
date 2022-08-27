package v1

type AddBannerRequest struct {
	BannerID int `json:"bannerId"`
	SlotID   int `json:"slotId"`
}

type RemoveBannerRequest struct {
	BannerID int `json:"bannerId"`
	SlotID   int `json:"slotId"`
}

type ClickBannerRequest struct {
	BannerID    int `json:"bannerId"`
	SlotID      int `json:"slotId"`
	UsergroupID int `json:"usergroupId"`
}
