package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	v1 "github.com/lozhkindm/otus-banners-rotation/internal/api/v1"
)

type Application interface {
	AddBanner(ctx context.Context, bannerID, slotID int) error
	RemoveBanner(ctx context.Context, bannerID, slotID int) error
	ClickBanner(ctx context.Context, bannerID, slotID, usergroupID int) error
	PickBanner(ctx context.Context, slotID, usergroupID int) (int, error)
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type Handlers struct {
	app    Application
	logger Logger
}

func NewHandlers(app Application, logger Logger) *Handlers {
	return &Handlers{
		app:    app,
		logger: logger,
	}
}

func (h *Handlers) AddBanner(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var request v1.AddBannerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := h.app.AddBanner(ctx, request.BannerID, request.SlotID); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) RemoveBanner(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var request v1.RemoveBannerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := h.app.RemoveBanner(ctx, request.BannerID, request.SlotID); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) ClickBanner(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var request v1.ClickBannerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := h.app.ClickBanner(ctx, request.BannerID, request.SlotID, request.UsergroupID); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) PickBanner(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sid := r.URL.Query().Get("slotId")
	uid := r.URL.Query().Get("usergroupId")
	if sid == "" || uid == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	slotID, err := strconv.Atoi(sid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	usergroupID, err := strconv.Atoi(uid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	bannerID, err := h.app.PickBanner(ctx, slotID, usergroupID)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := v1.PickBannerResponse{
		BannerID: bannerID,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
