package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/555ukr/bowatt/internal/websocket"
	"github.com/555ukr/bowatt/pkg/database"
	"github.com/555ukr/bowatt/pkg/models"
	"github.com/555ukr/bowatt/pkg/storage"
)

func UploadHandler(store storage.StorageService, repo database.PhotoRepository, hub *websocket.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("photo")
		if err != nil {
			http.Error(w, "missing 'photo' field", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "failed to read file", http.StatusInternalServerError)
			return
		}

		tagsRaw := r.FormValue("tags")
		var tags []string
		if tagsRaw != "" {
			for _, t := range strings.Split(tagsRaw, ",") {
				trimmed := strings.TrimSpace(t)
				if trimmed != "" {
					tags = append(tags, trimmed)
				}
			}
		}

		path, err := store.UploadFoto(header.Filename, fileBytes)
		if err != nil {
			http.Error(w, "failed to save file", http.StatusInternalServerError)
			return
		}

		photo := models.Photo{
			Path:      path,
			Tags:      tags,
			CreatedAt: time.Now(),
			Data:      base64.StdEncoding.EncodeToString(fileBytes),
		}

		if err := repo.InsertPhoto(r.Context(), photo); err != nil {
			log.Println("[ERROR]: ", err.Error())
			http.Error(w, "failed to save to database", http.StatusInternalServerError)
			return
		}

		err = hub.Broadcast(photo)
		if err != nil {
			// TODO: rollback insert and return error to the user
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(photo)
	}
}
