package songs

import "errors"

var (
	ErrDecodeRequest = errors.New("error decoding request")
	ErrUpdateSong    = errors.New("error updating song data")
	ErrDeleteSong    = errors.New("error deleting song data")
	ErrParseRequest  = errors.New("error parsing request")
	ErrGetSongs      = errors.New("error getting songs")
	ErrGetVerses     = errors.New("error getting verses")
	ErrMissingId     = errors.New("missing song id")
)
