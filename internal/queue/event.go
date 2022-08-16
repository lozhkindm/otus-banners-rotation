package queue

const (
	EventTypeImpression = iota
	EventTypeClick
)

type Event struct {
	Type        int   `json:"type"`
	BannerID    int   `json:"bannerId"`
	SlotID      int   `json:"slotId"`
	UsergroupID int   `json:"usergroupId"`
	CreatedAt   int64 `json:"createdAt"`
}
