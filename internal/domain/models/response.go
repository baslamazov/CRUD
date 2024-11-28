//go:generate easyjson -all response.go
package models

// easyjson:json
type SongResponse struct {
	ID          string  `json:"id"`
	GroupName   string  `json:"group_name"`
	Name        string  `json:"song_name"`
	ReleaseDate string  `json:"release_date"`
	Lyrics      []Lyric `json:"lyrics"`
	Link        string  `json:"link"`
}

// easyjson:json
type SongsResponse []SongResponse
