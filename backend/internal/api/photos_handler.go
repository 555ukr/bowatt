package api

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/555ukr/bowatt/pkg/database"
	"github.com/555ukr/bowatt/pkg/models"
	"github.com/555ukr/bowatt/pkg/storage"
)

type PhotosResponse struct {
	Photos     []models.Photo `json:"photos"`
	NextCursor string         `json:"next_cursor,omitempty"`
}

func ListPhotosHandler(store storage.StorageService, repo database.PhotoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tags []string
		tagsParam := r.URL.Query().Get("tags")
		if tagsParam != "" {
			for _, t := range strings.Split(tagsParam, ",") {
				trimmed := strings.TrimSpace(t)
				if trimmed != "" {
					tags = append(tags, trimmed)
				}
			}
		}

		limit := 20
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		var cursor *time.Time
		if c := r.URL.Query().Get("cursor"); c != "" {
			t, err := time.Parse(time.RFC3339Nano, c)
			if err != nil {
				http.Error(w, "invalid cursor format, use RFC3339", http.StatusBadRequest)
				return
			}
			cursor = &t
		}

		params := database.ListPhotosParams{
			Tags:   tags,
			Cursor: cursor,
			Limit:  limit,
		}

		photos, err := repo.ListPhotos(r.Context(), params)
		if err != nil {
			http.Error(w, "failed to query photos", http.StatusInternalServerError)
			return
		}

		for i := range photos {
			fileBytes, err := store.GetFoto(photos[i].Path)
			if err != nil {
				log.Printf("[ERROR]: failed to read file %s: %v", photos[i].Path, err)
				continue
			}
			photos[i].Data = base64.StdEncoding.EncodeToString(fileBytes)
		}

		resp := PhotosResponse{Photos: photos}
		if len(photos) == limit {
			lastPhoto := photos[len(photos)-1]
			resp.NextCursor = lastPhoto.CreatedAt.Format(time.RFC3339Nano)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
