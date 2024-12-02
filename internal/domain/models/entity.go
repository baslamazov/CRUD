package models

// Group
// easyjson:json
type Group struct {
	ID   string `json:"id" example:"1"`             // Уникальный идентификатор группы
	Name string `json:"name" example:"The Beatles"` // Название группы
}

// Song
// easyjson:json
type Song struct {
	ID          string `json:"id" example:"101"`                           // Уникальный идентификатор песни
	GroupName   string `json:"group_name" example:"The Beatles"`           // Название группы
	Name        string `json:"song" example:"Hey Jude"`                    // Название песни
	ReleaseDate string `json:"release_date" example:"1968-08-26"`          // Дата релиза
	Link        string `json:"link" example:"http://example.com/song/101"` // Ссылка на песню
}

// Lyric представляет текст песни
// easyjson:json
type Lyric struct {
	VerseNumber int    `json:"verse_number" example:"1"`                   // Номер куплета
	Text        string `json:"text" example:"Hey Jude, don't make it bad"` // Текст куплета
}

// SongDetail
// easyjson:json
type SongDetail struct {
	ReleaseDate string `json:"releaseDate" example:"1968-08-26"`           // Дата релиза
	Text        string `json:"text" example:"Hey Jude, don't make it bad"` // Текст песни
	Link        string `json:"link" example:"http://example.com/song/101"` // Ссылка на песню
}

// Songs
// easyjson:json
type Songs []Song

// Lyrics
// easyjson:json
type Lyrics []Lyric
