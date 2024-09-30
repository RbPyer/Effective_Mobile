package add

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"net/url"
	"rest/internal/lib/api/response"
	"rest/internal/server/handlers/songs"
	"rest/pkg/models"
	"time"
)

const ExternalAPIURL = "http://127.0.0.1:8080"

type SongSaver interface {
	AddSong(ctx context.Context, songDTO *models.SongDTO) (int64, error)
}

type CacheSaver interface {
	Set(ctx context.Context, songDTO *models.SongDTO) error
}

type Request struct {
	GroupName string `json:"group_name" validate:"required"`
	SongName  string `json:"song_name" validate:"required"`
}

type Response struct {
	response.Response
	Id int64 `json:"id"`
}

type ExternalAPIResponse struct {
	Text        string    `json:"text"`
	ReleaseDate time.Time `json:"release_date"`
	Link        string    `json:"link"`
	Error       string    `json:"error,omitempty"`
}

// @Summary Add a new song
// @Description Add a new song by group name and song name.
// @Tags songs
// @Accept json
// @Produce json
// @Param request body Request true "Song details"
// @Success 201 {object} Response "Successful add new song"
// @Failure 400 {object} response.Response "Bad request error response"
// @Failure 500 {object} response.Response "Internal server error response"
// @Router /songs [post]
func New(ctx context.Context, log *slog.Logger, songSaver SongSaver, cache CacheSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server/handlers/songs/add.New"

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

		log.Debug("request was successfully decoded", slog.String("op", op), slog.Any("request", req))

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

		groupName, songName := req.GroupName, req.SongName
		respFromExternalAPI, err := http.Get(
			fmt.Sprintf("%s/info?group_name=%s&song_name=%s",
				ExternalAPIURL,
				url.QueryEscape(groupName),
				url.QueryEscape(songName)),
		)

		if err != nil {
			log.Error("failed to send external API request",
				slog.String("op", op),
				slog.String("error", err.Error()),
			)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to send external API request"))
			return
		}
		defer func() {
			if err = respFromExternalAPI.Body.Close(); err != nil {
				log.Error("failed to close response body",
					slog.String("op", op),
					slog.String("error", err.Error()),
				)
			}
		}()

		var externalAPIResponse ExternalAPIResponse
		if err = render.DecodeJSON(respFromExternalAPI.Body, &externalAPIResponse); err != nil {
			log.Error("failed to parse external API request",
				slog.String("op", op),
				slog.String("error", err.Error()),
			)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to parse external API request, try again later"))
			return
		}

		if respFromExternalAPI.StatusCode != http.StatusOK {
			log.Error("error in external API response",
				slog.String("op", op),
				slog.Int("received status_code", respFromExternalAPI.StatusCode),
				slog.String("error", externalAPIResponse.Error),
			)
			render.Status(r, http.StatusBadGateway)
			render.JSON(w, r, response.Error(externalAPIResponse.Error))
			return
		}

		songDTO := models.SongDTO{
			GroupName:   groupName,
			SongName:    songName,
			ReleaseDate: externalAPIResponse.ReleaseDate,
			Text:        externalAPIResponse.Text,
			Link:        externalAPIResponse.Link,
		}
		id, err := songSaver.AddSong(ctx, &songDTO)
		if err != nil {
			log.Error("failed to add get", slog.String("op", op), slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to add new get"))
			return
		}

		if err = cache.Set(ctx, &songDTO); err != nil {
			log.Error("failed to set cache", slog.String("op", op), slog.String("error", err.Error()))
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{response.OK(), id})
		return
	}
}
