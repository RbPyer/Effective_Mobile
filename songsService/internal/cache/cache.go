package cache

import (
	"context"
	"fmt"
	"rest/pkg/models"
	"time"
)

type IDatabase interface {
	GetAllSongs(ctx context.Context) ([]models.SongDTO, error)
}

type ISongsLoader interface {
	Set(ctx context.Context, song *models.SongDTO) error
}

func Load(ctx context.Context, cache ISongsLoader, db IDatabase) error {
	const op = "cache.CacheLoad"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	songs, err := db.GetAllSongs(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, song := range songs {
		if err = cache.Set(ctx, &song); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil

}
