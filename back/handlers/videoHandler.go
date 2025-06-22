package handlers

import (
	"back/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type VideoHandler struct {
	Service services.VideoService
}

func NewVideoHandler(service services.VideoService) *VideoHandler {
	return &VideoHandler{Service: service}
}

func (h *VideoHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("user_id") // from JWT middleware

	title := r.FormValue("title")
	description := r.FormValue("description")

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File upload error: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = h.Service.UploadVideo(title, description, userID, file)
	if err != nil {
		http.Error(w, "Upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Video uploaded successfully"))
}

func (h *VideoHandler) GetVideos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	page := 1
	limit := 10
	if p := q.Get("page"); p != "" {
		if pi, err := strconv.Atoi(p); err == nil && pi > 0 {
			page = pi
		}
	}
	if l := q.Get("limit"); l != "" {
		if li, err := strconv.Atoi(l); err == nil && li > 0 {
			limit = li
		}
	}

	videos, err := h.Service.GetVideosPaginated(page, limit)
	if err != nil {
		http.Error(w, "Error fetching videos", http.StatusInternalServerError)
		return
	}

	type VideoResponse struct {
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		URL          string    `json:"url"`
		OptimizedURL string    `json:"optimized_url,omitempty"`
		UploaderID   string    `json:"uploader_id"`
		UploadDate   time.Time `json:"upload_date"`
	}
	var resp []VideoResponse
	for _, v := range videos {
		optURL := ""
		if v.URL != "" {

			optURL = strings.Replace(v.URL, "/upload/", "/upload/q_auto,f_auto/", 1)
		}
		resp = append(resp, VideoResponse{
			Title:        v.Title,
			Description:  v.Description,
			URL:          v.URL,
			OptimizedURL: optURL,
			UploaderID:   v.UploaderID,
			UploadDate:   v.UploadDate,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
