package verses

import (
	"context"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"rest/internal/lib/api/response"
	"rest/internal/server/handlers/songs"
	"strconv"
)

type VerseGetter interface {
	GetVerses(ctx context.Context, id int64, offset, limit int) (string, error)
}

type Response struct {
	response.Response
	Verses string `json:"verses"`
}

// @Summary Get song verses
// @Description Retrieve verses of a song by its ID with optional pagination.
// @Tags songs
// @Accept json
// @Produce json
// @Param id query int true "Song ID"
// @Param verse query int false "Verse number" default(1)
// @Param limit query int false "Number of verses to retrieve" default(1)
// @Success 200 {object} Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /verses [get]
func New(ctx context.Context, log *slog.Logger, verseGetter VerseGetter, cache VerseGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server/handlers/songs/verses.New"
		log.Debug("Start handling...", slog.String("op", op))

		verse, limit := 1, 1

		err := r.ParseForm()
		if err != nil {
			log.Error("unable to parse form", slog.String("op", op), slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(songs.ErrParseRequest.Error()))
			return
		}

		if r.FormValue("id") == "" {
			log.Error("unable to parse required id field", slog.String("op", op))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(songs.ErrMissingId.Error()))
			return
		}

		id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if err != nil {
			log.Error("unable to parse id", slog.String("op", op), slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(songs.ErrParseRequest.Error()))
			return
		}

		if r.FormValue("verse") != "" {
			verse, err = strconv.Atoi(r.FormValue("verse"))
			if err != nil {
				log.Error("unable to parse verse",
					slog.String("op", op),
					slog.String("error", err.Error()),
				)
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error(songs.ErrParseRequest.Error()))
				return
			}
		}

		if r.FormValue("limit") != "" {
			limit, err = strconv.Atoi(r.FormValue("limit"))
			if err != nil {
				log.Error("unable to parse limit",
					slog.String("op", op),
					slog.String("error", err.Error()),
				)
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error(songs.ErrParseRequest.Error()))
				return
			}
		}

		if verse < 1 || limit < 1 {
			log.Error("verse and limit must be greater than 0", slog.String("op", op))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("verse and limit must be greater than 0"))
			return
		}

		resultVerses, err := cache.GetVerses(ctx, id, verse, limit)
		if err != nil {
			log.Error("unable to get verses from cache",
				slog.String("op", op),
				slog.String("error", err.Error()),
			)

			resultVerses, err = verseGetter.GetVerses(ctx, id, verse, limit)
			if err != nil {
				log.Error("unable to get verses",
					slog.String("op", op),
					slog.String("error", err.Error()),
				)
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, response.Error(songs.ErrGetVerses.Error()))
				return
			}

			log.Debug("verses were received", slog.String("op", op))

			render.Status(r, http.StatusOK)
			render.JSON(w, r, Response{
				response.OK(),
				resultVerses,
			})
			return
		}
		log.Debug("verses from cache were received", slog.String("op", op))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, Response{
			response.OK(),
			resultVerses,
		})
		return
	}
}
