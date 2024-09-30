package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"rest/internal/storage"
	"rest/pkg/models"
	"strings"
	"time"
)

const (
	setOP   = "SET"
	whereOP = "WHERE"
)

type Storage struct {
	dbPool *pgxpool.Pool
}

func (s *Storage) GetPoolForGracefulShutdown() *pgxpool.Pool {
	return s.dbPool
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	const op = "storage/postgres.New"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{dbPool: dbPool}, nil
}

func (s *Storage) AddSong(ctx context.Context, songDTO *models.SongDTO) (int64, error) {
	const op = "storage/postgres.AddSong"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := s.dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)

	}

	const query = `INSERT INTO songs(group_name, song_name, release_date, song_text, link) 
                   VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	var id int64
	err = tx.QueryRow(ctx, query, songDTO.GroupName, songDTO.SongName, songDTO.ReleaseDate,
		songDTO.Text, songDTO.Link).Scan(&id)
	if err != nil {
		transactionRollback(ctx, tx, op)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		transactionRollback(ctx, tx, op)
		return 0, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetAllSongs(ctx context.Context) ([]models.SongDTO, error) {
	const op = "storage/postgres.GetAllSongs"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, group_name, song_name, release_date, song_text, link FROM songs`

	rows, err := s.dbPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	var songs []models.SongDTO
	for rows.Next() {
		var song models.SongDTO
		if err = rows.Scan(&song.Id, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (s *Storage) GetSongs(ctx context.Context, songDTO *models.SongDTO, page int, limit int,
	opts ...models.OptionFunc) ([]models.SongDTO, error) {

	const op = "storage/postgres.GetSongs"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	offset := (page - 1) * limit

	query := `SELECT id, group_name, song_name, release_date, song_text, link FROM songs `
	whereClause, args := models.BuildQuery(whereOP, songDTO, opts...)
	query += whereClause
	query += fmt.Sprintf(" ORDER BY id OFFSET $%d LIMIT $%d", len(args)+1, len(args)+2)
	args = append(args, offset, limit)

	rows, err := s.dbPool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	songs := make([]models.SongDTO, 0)
	for rows.Next() {
		var song models.SongDTO
		if err = rows.Scan(&song.Id, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (s *Storage) GetVerses(ctx context.Context, id int64, verse, limit int) (string, error) {
	const op = "storage/postgres.GetVerses"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT song_text FROM songs WHERE id=$1`
	var songText string
	err := s.dbPool.QueryRow(ctx, query, id).Scan(&songText)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	verses := strings.Split(songText, "\n\n")

	start := (verse - 1) * limit
	end := start + limit
	if start >= len(verses) {
		start = 0
	}
	if end > len(verses) {
		end = 1
	}

	return strings.Join(verses[start:end], "\n\n"), nil
}

func (s *Storage) UpdateSong(ctx context.Context, songDTO *models.SongDTO, opts ...models.OptionFunc) error {
	const op = "storage/postgres.UpdateSong"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := s.dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := `UPDATE songs `
	setClause, args := models.BuildQuery(setOP, songDTO, opts...)
	query += setClause
	query += fmt.Sprintf(" WHERE id=$%d", len(args)+1)
	args = append(args, songDTO.Id)

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		transactionRollback(ctx, tx, op)
		return fmt.Errorf("%s: %w", op, err)
	}
	if err = tx.Commit(ctx); err != nil {
		transactionRollback(ctx, tx, op)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteSong(ctx context.Context, id int64) error {
	const op = "storage/postgres.DeleteSong"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := s.dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := tx.Exec(ctx, "DELETE FROM songs WHERE id = $1;", id)
	if err != nil {
		transactionRollback(ctx, tx, op)
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.RowsAffected() == 0 {
		transactionRollback(ctx, tx, op)
		return storage.ErrNoAffected
	}

	if err = tx.Commit(ctx); err != nil {
		transactionRollback(ctx, tx, op)
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func transactionRollback(ctx context.Context, tx pgx.Tx, op string) {
	err := tx.Rollback(ctx)
	if err != nil {
		panicMessage := fmt.Sprintf("%s: %s", op, err.Error())
		panic(panicMessage)
	}
}
