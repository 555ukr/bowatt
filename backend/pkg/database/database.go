package database

import (
	"context"
	"database/sql"

	"github.com/555ukr/bowatt/pkg/models"
	"github.com/lib/pq"
)

// PhotoRepository defines the interface for photo persistence.
type PhotoRepository interface {
	InsertPhoto(ctx context.Context, photo models.Photo) error
}

// PostgresPhotoRepository implements PhotoRepository using PostgreSQL.
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
