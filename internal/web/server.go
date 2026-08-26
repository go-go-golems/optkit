package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-go-golems/optkit/query"
	"github.com/go-go-golems/optkit/record"
)

//go:embed static/index.html static/styles.css static/app.js
var staticFiles embed.FS

type Server struct {
	Query        query.Service
	PollInterval time.Duration
}

func (s Server) Handler() (http.Handler, error) {
	if s.Query.Metadata == nil || s.Query.Artifacts == nil {
		return nil, fmt.Errorf("web server requires a query service")
	}
	assets, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", s.handleHealth)
	mux.HandleFunc("GET /api/v1/campaigns", s.handleCampaigns)
	mux.HandleFunc("GET /api/v1/campaigns/{campaignID}", s.handleCampaign)
	mux.HandleFunc("GET /api/v1/campaigns/{campaignID}/events", s.handleEvents)
	mux.HandleFunc("GET /api/v1/campaigns/{campaignID}/stream", s.handleStream)
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(assets))))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		data, readErr := fs.ReadFile(assets, "index.html")
		if readErr != nil {
			http.Error(w, "web application unavailable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(data)
	})
	return securityHeaders(mux), nil
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; connect-src 'self'")
		next.ServeHTTP(w, r)
	})
}

func (s Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "api_version": query.APIVersion, "read_only": true})
}

func (s Server) handleCampaigns(w http.ResponseWriter, r *http.Request) {
	page, err := s.Query.ListCampaigns(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "campaign_list_failed", err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (s Server) handleCampaign(w http.ResponseWriter, r *http.Request) {
	id, ok := campaignID(w, r)
	if !ok {
		return
	}
	view, err := s.Query.Campaign(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "campaign_not_found", err)
		return
	}
	w.Header().Set("ETag", fmt.Sprintf("%q", fmt.Sprintf("%s:%d:%s", id, view.Version, view.Overview.LastDigest)))
	writeJSON(w, http.StatusOK, view)
}

func (s Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	id, ok := campaignID(w, r)
	if !ok {
		return
	}
	after, err := uintQuery(r, "after", 0)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_after", err)
		return
	}
	limit, err := intQuery(r, "limit", query.DefaultEventLimit)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_limit", err)
		return
	}
	page, err := s.Query.Events(r.Context(), id, after, limit)
	if err != nil {
		writeError(w, http.StatusNotFound, "events_unavailable", err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (s Server) handleStream(w http.ResponseWriter, r *http.Request) {
	id, ok := campaignID(w, r)
	if !ok {
		return
	}
	cursor, err := uintQuery(r, "after", 0)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_after", err)
		return
	}
	if raw := strings.TrimSpace(r.Header.Get("Last-Event-ID")); raw != "" {
		parsed, parseErr := strconv.ParseUint(raw, 10, 64)
		if parseErr != nil {
			writeError(w, http.StatusBadRequest, "invalid_last_event_id", parseErr)
			return
		}
		if parsed > cursor {
			cursor = parsed
		}
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "stream_unsupported", fmt.Errorf("response writer cannot flush"))
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	interval := s.PollInterval
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		page, pageErr := s.Query.Events(r.Context(), id, cursor, query.DefaultEventLimit)
		if pageErr != nil {
			writeSSEError(w, pageErr)
			flusher.Flush()
			return
		}
		for _, event := range page.Events {
			data, marshalErr := json.Marshal(event)
			if marshalErr != nil {
				writeSSEError(w, marshalErr)
				flusher.Flush()
				return
			}
			_, _ = fmt.Fprintf(w, "id: %d\nevent: campaign-event\ndata: %s\n\n", event.Seq, data)
			cursor = event.Seq
		}
		if len(page.Events) > 0 {
			flusher.Flush()
			continue
		}
		_, _ = fmt.Fprint(w, ": heartbeat\n\n")
		flusher.Flush()
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
		}
	}
}

func campaignID(w http.ResponseWriter, r *http.Request) (record.CampaignID, bool) {
	raw := r.PathValue("campaignID")
	if err := record.ValidateID("campaign", raw); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_campaign_id", err)
		return "", false
	}
	return record.CampaignID(raw), true
}

func uintQuery(r *http.Request, name string, fallback uint64) (uint64, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(name))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a non-negative integer", name)
	}
	return value, nil
}

func intQuery(r *http.Request, name string, fallback int) (int, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(name))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer", name)
	}
	return value, nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, code string, err error) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": err.Error()}})
}

func writeSSEError(w http.ResponseWriter, err error) {
	data, _ := json.Marshal(map[string]string{"code": "stream_failed", "message": err.Error()})
	_, _ = fmt.Fprintf(w, "event: error\ndata: %s\n\n", data)
}

var _ interface {
	Handler() (http.Handler, error)
} = Server{}
