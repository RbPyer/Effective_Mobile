package update

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"rest/internal/lib/api/response"
	tp "rest/internal/lib/timeParser"
	"rest/internal/server/handlers/songs"
	"rest/pkg/models"
)

type SongUpdater interface {
	UpdateSong(ctx context.Context, songDTO *models.SongDTO, opts ...models.OptionFunc) error
}

type CacheUpdater interface {
	Update(ctx context.Context, song *models.SongDTO) error
}

type Request struct {
	Id          int64  `json:"id" validate:"required"`
	GroupName   string `json:"group_name,omitempty"`
	SongName    string `json:"song_name,omitempty"`
	Text        string `json:"text,omitempty"`
	ReleaseDate string `json:"release_date,omitempty"`
	Link        string `json:"link,omitempty" validate:"omitempty,url"`
}

// @Summary Update a song
// @Description Update a song's details by its ID.
// @Tags songs
// @Accept json
// @Produce json
// @Param request body Request true "Updated song details"
// @Success 204 {object} nil
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /songs [put]
func New(ctx context.Context, log *slog.Logger, songUpdater SongUpdater, cd CacheUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server/handlers/songs/update.New"

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error(songs.ErrDecodeRequest.Error(),
				slog.String("op", op),
				slog.String("error", err.Error()),
			)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(songs.ErrDecodeRequest.Error()))
			return
		}

		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Error("Failed to close request body: %v",
					slog.String("op", op),
					slog.String("error", err.Error()),
				)
			}
		}()

		log.Debug("request was successfully decoded",
			slog.String("op", op),
			slog.Any("request", req),
		)

		if err := validator.New().Struct(req); err != nil {
			var validationErrs validator.ValidationErrors
			errors.As(err, &validationErrs)
			log.Error("failed to validate request",
				slog.String("op", op),
				slog.String("error", err.Error()),
			)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validationErrs))
			return
		}

		date, err := tp.ParseToDate(req.ReleaseDate)
		if err != nil {
			log.Error("failed to parse release date",
				slog.String("op", op),
				slog.String("error", err.Error()),
			)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid release date"))
			return
		}

		songDTO := &models.SongDTO{
			Id:          req.Id,
			GroupName:   req.GroupName,
			SongName:    req.SongName,
			Text:        req.Text,
			ReleaseDate: date,
			Link:        req.Link,
		}

		err = songUpdater.UpdateSong(
			ctx,
			songDTO,
			models.WithGroupName(),
			models.WithSongName(),
			models.WithReleaseDate(),
			models.WithLink(),
			models.WithSongText(),
		)

		if err != nil {
			log.Error(songs.ErrUpdateSong.Error(),
				slog.String("op", op),
				slog.String("error", err.Error()),
			)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(songs.ErrUpdateSong.Error()))
			return
		}

		if err = cd.Update(ctx, songDTO); err != nil {
			log.Error("failed to update song data in cache",
				slog.String("op", op),
				slog.String("error", err.Error()),
			)
		}

		render.Status(r, http.StatusNoContent)
		render.JSON(w, r, nil)
		return
	}
}
