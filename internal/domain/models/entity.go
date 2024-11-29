//go:generate easyjson -all entity.go

package models

//easyjson:json
type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//easyjson:json
type Song struct {
	ID          string `json:"id"`
	GroupName   string `json:"group_name"`
	Name        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Link        string `json:"link"`
}

//easyjson:json
type Lyric struct {
	VerseNumber int    `json:"verse_number"`
	Text        string `json:"text"`
}

//easyjson:json
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

//easyjson:json
type Songs []Song

//easyjson:json
type Lyrics []Lyric
