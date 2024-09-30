package info

import (
	"fmt"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

type Response struct {
	Text        string    `json:"text"`
	ReleaseDate time.Time `json:"release_date"`
	Link        string    `json:"link"`
	Error       string    `json:"error,omitempty"`
}

// @Summary Get song info
// @Description Get information about a song based on group name and song name.
// @Tags info
// @Accept json
// @Produce json
// @Param group_name query string true "Name of the music group"
// @Param song_name query string true "Name of the song"
// @Success 200 {object} Response "Successful response with song info"
// @Failure 400 {object} Response "Bad request error response"
// @Router /info [get]
func New(counter *atomic.Uint32, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/info.New"

		groupName := r.URL.Query().Get("group_name")
		if groupName == "" {
			log.Error("group_name was not received",
				slog.String("op", op),
				slog.String("url", r.URL.String()),
			)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{Error: "group_name is required"})
			return
		}
		songName := r.URL.Query().Get("song_name")
		if songName == "" {
			log.Error("song_name was not received",
				slog.String("op", op),
				slog.String("url", r.URL.String()),
			)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{Error: "song_name is required"})
			return
		}

		now := time.Now()
		releaseDate := time.Date(now.Year(),
			now.Month(),
			now.Day()-int(counter.Load()),
			0,
			0,
			0,
			0,
			time.UTC,
		)

		groupName = strings.ReplaceAll(groupName, " ", "_")
		songName = strings.ReplaceAll(groupName, " ", "_")

		resp := Response{
			Text:        "verse1\n\nverse2\n\nverse3",
			Link:        fmt.Sprintf("https://songs.ru/%s/%s", groupName, groupName),
			ReleaseDate: releaseDate,
		}

		counter.Add(1)
		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
		return
	}
}
