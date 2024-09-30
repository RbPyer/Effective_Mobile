package redis

import (
	"context"
	"fmt"
	r "github.com/go-redis/redis"
	"rest/pkg/models"
	"strconv"
	"strings"
	"time"
)

type Cache struct {
	redisClient *r.Client
}

func New(ctx context.Context, addr string) (*Cache, error) {
	const op = "cache/redis.New"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rcl := r.NewClient(&r.Options{
		Addr: addr,
	})
	err := rcl.Ping().Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Cache{redisClient: rcl}, nil
}

func (c *Cache) Set(ctx context.Context, song *models.SongDTO) error {
	const op = "cache/redis.Set"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := map[string]interface{}{
		"group_name":   song.GroupName,
		"song_name":    song.SongName,
		"release_date": song.ReleaseDate,
		"link":         song.Link,
		"song_text":    song.Text,
	}

	err := c.redisClient.HMSet(strconv.FormatInt(song.Id, 10), m).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Cache) Update(ctx context.Context, song *models.SongDTO) error {
	const op = "cache/redis.Update"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := make(map[string]interface{}, 5)
	if song.GroupName != "" {
		m["group_name"] = song.GroupName
	}
	if song.SongName != "" {
		m["song_name"] = song.SongName
	}
	if song.Link != "" {
		m["link"] = song.Link
	}
	if song.Text != "" {
		m["song_text"] = song.Text
	}
	if !song.ReleaseDate.IsZero() {
		m["release_date"] = song.ReleaseDate
	}

	err := c.redisClient.HMSet(strconv.FormatInt(song.Id, 10), m).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Cache) Del(ctx context.Context, id int64) error {
	const op = "cache/redis.Del"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := c.redisClient.HDel(strconv.FormatInt(id, 10)).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Cache) GetVerses(ctx context.Context, id int64, verse, limit int) (string, error) {
	const op = "cache/redis.GetVerse"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	songText := c.redisClient.HGet(strconv.FormatInt(id, 10), "song_text").Val()
	if songText == "" {
		return "", fmt.Errorf("%s: %w", op, fmt.Errorf("no song found with id %d", id))
	}
	c.redisClient.Keys("song_")

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
