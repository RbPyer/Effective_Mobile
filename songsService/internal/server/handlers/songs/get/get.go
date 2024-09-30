package get

import (
	"context"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"rest/internal/lib/api/response"
	tp "rest/internal/lib/timeParser"
	"rest/internal/server/handlers/songs"
	"rest/pkg/models"
	"strconv"
	"time"
)

type SongsGetter interface {
	GetSongs(ctx context.Context, songDTO *models.SongDTO, page int, limit int, opts ...models.OptionFunc) ([]models.SongDTO, error)
}

type Response struct {
	response.Response
	Songs []models.SongDTO `json:"songs"`
}

// @Summary Get songs
// @Description Retrieve songs with pagination and filtering options.
// @Tags songs
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of results per page" default(10)
// @Param id query int false "Song ID"
// @Param release_date query string false "Release date"
// @Param group_name query string false "Group name"
// @Param song_name query string false "Song name"
// @Param song_text query string false "Song text"
// @Param link query string false "Song link"
// @Success 200 {object} Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /songs [get]
func New(ctx context.Context, log *slog.Logger, songsGetter SongsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server/handlers/songs/get.New"

		page, limit := 1, 10

		songDTO := &models.SongDTO{}

		err := r.ParseForm()
		if err != nil {
			log.Error("unable to parse form", slog.String("op", op), slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(songs.ErrParseRequest.Error()))
			return
		}

		if r.FormValue("page") != "" {
			page, err = strconv.Atoi(r.FormValue("page"))
			if err != nil {
				log.Error("unable to parse page",
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

		if r.FormValue("id") != "" {
			var id int64
			id, err = strconv.ParseInt(r.FormValue("id"), 10, 64)
			if err != nil {
				log.Error("unable to parse id",
					slog.String("op", op),
					slog.String("error", err.Error()),
				)
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error(songs.ErrParseRequest.Error()))
				return
			}
			songDTO.Id = id
		}

		if r.FormValue("release_date") != "" {
			var releaseDate time.Time
			releaseDate, err = tp.ParseToDate(r.FormValue("release_date"))
			if err != nil {
				log.Error("unable to parse release_date",
					slog.String("op", op),
					slog.String("error", err.Error()),
				)
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error(songs.ErrParseRequest.Error()))
				return
			}
			songDTO.ReleaseDate = releaseDate
		}

		songDTO.GroupName = r.FormValue("group_name")
		songDTO.SongName = r.FormValue("song_name")
		songDTO.Text = r.FormValue("song_text")
		songDTO.Link = r.FormValue("link")

		songsData, err := songsGetter.GetSongs(
			ctx,
			songDTO,
			page,
			limit,
			models.WithId(),
			models.WithReleaseDate(),
			models.WithGroupName(),
			models.WithSongName(),
			models.WithLink(),
			models.WithSongText(),
		)
		if err != nil {
			log.Error("unable to get songs", slog.String("op", op), slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(songs.ErrGetSongs.Error()))
			return
		}

		log.Debug("songs were received", slog.String("op", op))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, Response{
			response.OK(),
			songsData,
		})
		return
	}
}
