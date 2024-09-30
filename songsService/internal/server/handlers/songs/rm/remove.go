package rm

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"rest/internal/lib/api/response"
	"rest/internal/server/handlers/songs"
	"strconv"
)

type SongDeleter interface {
	DeleteSong(ctx context.Context, id int64) error
}

type CacheDeleter interface {
	Del(ctx context.Context, id int64) error
}

// @Summary Remove a song
// @Description Delete a song by its ID.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 204 {object} nil
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /songs/{id} [delete]
func New(ctx context.Context, log *slog.Logger, songDeleter SongDeleter, cd CacheDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server/handlers/songs/rm.New"

		id := chi.URLParam(r, "id")
		if id == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("missing id parameter"))
			return
		}

		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Error("Failed to parse id parameter: %v",
				slog.String("op", op),
				slog.String("error", err.Error()),
			)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid id parameter"))
			return
		}

		err = songDeleter.DeleteSong(ctx, idInt)
		if err != nil {
			log.Error(songs.ErrDeleteSong.Error(), slog.String("op", op), slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(songs.ErrDeleteSong.Error()))
			return
		}

		err = cd.Del(r.Context(), idInt)
		if err != nil {
			log.Error("Failed to remove song from cache",
				slog.String("op", op),
				slog.String("error", err.Error()),
			)
		}

		render.Status(r, http.StatusNoContent)
		render.JSON(w, r, nil)
		return
	}
}
