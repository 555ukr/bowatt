package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/555ukr/bowatt/pkg/models"
	"github.com/lib/pq"
)

type PhotoRepository interface {
	InsertPhoto(ctx context.Context, photo models.Photo) error
	ListPhotos(ctx context.Context, params ListPhotosParams) ([]models.Photo, error)
}

type ListPhotosParams struct {
	Tags   []string
	Cursor *time.Time
	Limit  int
}

type PostgresPhotoRepository struct {
	DB *sql.DB
}

func NewPostgresPhotoRepository(db *sql.DB) PhotoRepository {
	return &PostgresPhotoRepository{DB: db}
}

func (r *PostgresPhotoRepository) InsertPhoto(ctx context.Context, photo models.Photo) error {
	query := `INSERT INTO photo (path, tags, created_at) VALUES ($1, $2, $3)`
	_, err := r.DB.ExecContext(ctx, query, photo.Path, pq.Array(photo.Tags), photo.CreatedAt)
	return err
}

func (r *PostgresPhotoRepository) ListPhotos(ctx context.Context, params ListPhotosParams) ([]models.Photo, error) {
	var rows *sql.Rows
	var err error

	limit := params.Limit
	if limit <= 0 {
		limit = 20
	}

	hasTags := len(params.Tags) > 0
	hasCursor := params.Cursor != nil

	switch {
	case hasTags && hasCursor:
		query := `SELECT path, tags, created_at FROM photo WHERE tags && $1 AND created_at < $2 ORDER BY created_at DESC LIMIT $3`
		rows, err = r.DB.QueryContext(ctx, query, pq.Array(params.Tags), *params.Cursor, limit)
	case hasTags:
		query := `SELECT path, tags, created_at FROM photo WHERE tags && $1 ORDER BY created_at DESC LIMIT $2`
		rows, err = r.DB.QueryContext(ctx, query, pq.Array(params.Tags), limit)
	case hasCursor:
		query := `SELECT path, tags, created_at FROM photo WHERE created_at < $1 ORDER BY created_at DESC LIMIT $2`
		rows, err = r.DB.QueryContext(ctx, query, *params.Cursor, limit)
	default:
		query := `SELECT path, tags, created_at FROM photo ORDER BY created_at DESC LIMIT $1`
		rows, err = r.DB.QueryContext(ctx, query, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []models.Photo
	for rows.Next() {
		var p models.Photo
		if err := rows.Scan(&p.Path, pq.Array(&p.Tags), &p.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}

	return photos, rows.Err()
}
