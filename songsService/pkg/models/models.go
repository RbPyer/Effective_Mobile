package models

import "time"

type SongDTO struct {
	Id          int64
	GroupName   string
	SongName    string
	Text        string
	ReleaseDate time.Time
	Link        string
}

func checkEmpty(s SongDTO) bool {
	return s.Id == 0 && s.GroupName == "" && s.SongName == "" && s.Text == "" && s.Link == "" && s.ReleaseDate.IsZero()
}
