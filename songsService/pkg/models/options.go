/*
Due to the small number of endpoints in this project, it was decided not to include query builders like Squirrel or ORMs like GORM.
Instead, a custom solution was developed, utilizing handwritten optional parameters combined with a build function.
This approach provides flexibility and simplicity without the overhead of external libraries, ensuring that the codebase remains lightweight while still being maintainable.
*/

package models

import (
	"fmt"
	"strings"
)

type OptionFunc func(songDTO *SongDTO, args []interface{}) (string, []interface{})

func BuildQuery(filter string, songDTO *SongDTO, opts ...OptionFunc) (string, []interface{}) {
	if checkEmpty(*songDTO) {
		return "", []interface{}{}
	}

	var clause string

	switch filter {
	case "SET":
		clause = "SET "
	case "WHERE":
		clause = "WHERE "
	}

	var row string
	updates := make([]string, 0, 6)
	args := make([]interface{}, 0, 6)

	for _, opt := range opts {
		row, args = opt(songDTO, args)
		if row != "" {
			updates = append(updates, row)
		}
	}

	clause += strings.Join(updates, ", ")

	return clause, args
}

// WithId is option-func for id field.
// Do not use in update-queries, use this opt only in select-queries, otherwise error will be occurred.
func WithId() OptionFunc {
	return func(songDTO *SongDTO, args []interface{}) (string, []interface{}) {
		if songDTO.Id != 0 {
			args = append(args, songDTO.Id)
			return fmt.Sprintf("id=$%d", len(args)), args
		}
		return "", args
	}
}

// WithGroupName is option-func for group_name field.
func WithGroupName() OptionFunc {
	return func(songDTO *SongDTO, args []interface{}) (string, []interface{}) {
		if songDTO.GroupName != "" {
			args = append(args, songDTO.GroupName)
			return fmt.Sprintf("group_name=$%d", len(args)), args
		}
		return "", args
	}
}

// WithSongName is option-func for song_name field.
func WithSongName() OptionFunc {
	return func(songDTO *SongDTO, args []interface{}) (string, []interface{}) {
		if songDTO.SongName != "" {
			args = append(args, songDTO.SongName)
			return fmt.Sprintf("song_name=$%d", len(args)), args
		}
		return "", args
	}
}

// WithReleaseDate is option-func for release_date field.
func WithReleaseDate() OptionFunc {
	return func(songDTO *SongDTO, args []interface{}) (string, []interface{}) {
		if !songDTO.ReleaseDate.IsZero() {
			args = append(args, songDTO.ReleaseDate)
			return fmt.Sprintf("release_date=$%d", len(args)), args
		}
		return "", args
	}
}

// WithLink is option-func for link field.
func WithLink() OptionFunc {
	return func(songDTO *SongDTO, args []interface{}) (string, []interface{}) {
		if songDTO.Link != "" {
			args = append(args, songDTO.Link)
			return fmt.Sprintf("link=$%d", len(args)), args
		}
		return "", args
	}
}

// WithSongText is option-func for song_text field.
func WithSongText() OptionFunc {
	return func(songDTO *SongDTO, args []interface{}) (string, []interface{}) {
		if songDTO.Text != "" {
			args = append(args, songDTO.Text)
			return fmt.Sprintf("song_text=$%d", len(args)), args
		}
		return "", args
	}
}
