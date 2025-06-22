package handlers

import (
	"back/services"
	"encoding/json"
	"net/http"
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

	videos, err := h.Service.GetAllVideos()
	if err != nil {
		http.Error(w, "Error fetching videos", http.StatusInternalServerError)
		return
	}

	type VideoResponse struct {
		Title       string    `json:"title"`
		URL         string    `json:"url"`
		Description string    `json:"description"`
		UploaderID  string    `json:"uploader_id"`
		UploadDate  time.Time `json:"upload_date"`
	}

	var resp []VideoResponse
	for _, v := range videos {
		resp = append(resp, VideoResponse{
			Title:       v.Title,
			URL:         v.URL,
			Description: v.Description,
			UploaderID:  v.UploaderID,
			UploadDate:  v.UploadDate,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
