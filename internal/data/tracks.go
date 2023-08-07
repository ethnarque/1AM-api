package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pmlogist/1am-audio-srv/internal/validator"
)

type Track struct {
	ID        int64         `json:"id"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updated_at"`
	Title     string        `json:"title"`
	Duration  time.Duration `json:"duration"`
	Genres    []string      `json:"genres,omitempty"`
	Albums    []string      `json:"albums,omitempty"`
	Version   uuid.UUID     `json:"version"`
}

type TrackModel struct {
	DB *pgxpool.Pool
}

func (m TrackModel) Create(track *Track) error {
	query := `
        INSERT INTO tracks (title, duration, genres, albums)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at, version
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dest := []any{&track.ID, &track.CreatedAt, &track.UpdatedAt, &track.Version}

	return m.DB.QueryRow(
		ctx,
		query,
		track.Title,
		track.Duration,
		track.Genres,
		track.Albums,
	).Scan(dest...)
}

func (m TrackModel) FindAll(title string, genres []string, filters Filter) ([]*Track, Metadata, error) {
	query := fmt.Sprintf(`
        SELECT count(*) OVER(), id, created_at, updated_at, title, duration, genres, albums, version
        FROM tracks
        WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
        AND (genres @> $2 OR $2 = '{}')
        ORDER BY %s %s, id ASC
        LIMIT $3 OFFSET $4
        `, filters.sortColumn(), filters.sortDirector())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, genres, filters.limit(), filters.offset()}

	rows, err := m.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	totalRecords := 0
	tracks := []*Track{}

	for rows.Next() {
		var track Track

		err := rows.Scan(
			&totalRecords,
			&track.ID,
			&track.CreatedAt,
			&track.UpdatedAt,
			&track.Title,
			&track.Duration,
			&track.Genres,
			&track.Albums,
			&track.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		tracks = append(tracks, &track)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tracks, metadata, nil
}

func (m TrackModel) Find(id int64) (*Track, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, created_at, updated_at, title, duration, genres, albums, version
        FROM tracks
        WHERE id = $1
    `

	var track Track

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, id).Scan(
		&track.ID,
		&track.CreatedAt,
		&track.UpdatedAt,
		&track.Title,
		&track.Duration,
		&track.Genres,
		&track.Albums,
		&track.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound

		default:
			return nil, err
		}
	}
	return &track, nil
}

func (m TrackModel) Update(track *Track) error {
	query := `
        UPDATE tracks
        SET title = $1, duration = $2, genres = $3, albums = $4, version = uuid_generate_v4()
        WHERE id = $5 AND version = $6
        RETURNING updated_at, version
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(
		ctx,
		query,
		track.Title,
		track.Duration,
		track.Genres,
		track.Albums,
		track.ID,
		track.Version,
	).Scan(&track.UpdatedAt, &track.Version)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m TrackModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
        DELETE FROM tracks
        WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if r := result.RowsAffected(); r == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func ValidateTrack(v *validator.Validator, track *Track) {
	v.Check(track.Title != "", "title", "must be provided")
	v.Check(len(track.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(track.Duration != 0, "duration", "must be provided")
	v.Check(track.Genres != nil, "genres", "must be provided")
	v.Check(len(track.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(track.Albums != nil, "albums", "must be provided")
	v.Check(len(track.Albums) >= 1, "albums", "must contain at least 1 genre")
}
