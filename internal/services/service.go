package services

import (
	"EffectiveMobile/internal/domain/models"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Service struct {
	log             *slog.Logger
	songDeleter     SongDeleter
	songSaver       SongSaver
	libraryProvider LibraryProvider
	tokenTTL        time.Duration
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SongSaver
type SongSaver interface {
	InsertSong(
		ctx context.Context,
		song models.Song,
		songDetail models.SongDetail,
	) (bool, error)
	UpdateSong(
		ctx context.Context,
		song models.Song,
	) (bool, error)
}
type SongDeleter interface {
	DeleteSong(
		ctx context.Context,
		song string,
		group string,
	) (bool, error)
}

type LibraryProvider interface {
	Song(ctx context.Context,
		songName string,
		groupName string,
		releaseDate string,
		page int,
		limit int,
	) (models.Songs, error)
	Lyric(ctx context.Context,
		songName string,
		groupName string,
		page int,
		limit int,
	) (models.Lyrics, error)
}

func New(
	log *slog.Logger,
	songSaver SongSaver,
	songDeleter SongDeleter,
	libraryProvider LibraryProvider,
) *Service {
	return &Service{
		songSaver:       songSaver,
		songDeleter:     songDeleter,
		libraryProvider: libraryProvider,
		log:             log,
	}
}
func (s *Service) NewSong(ctx context.Context,
	song models.Song,
) (bool, error) {

	apiUrl := os.Getenv("API_URL")

	// Проверка наличия слеша в конце apiUrl
	if !strings.HasSuffix(apiUrl, "/") {
		apiUrl += "/"
	}

	// Формирование параметров запроса
	params := url.Values{}
	params.Add("group", song.GroupName)
	params.Add("song", song.Name)

	// Формирование полного URL запроса
	apiEndpoint := apiUrl + "info" + "?" + params.Encode()

	// Запрос к внешнему API
	resp, err := http.Get(apiEndpoint)
	if err != nil {
		s.log.Error("Ошибка запроса к внешнему API")
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.log.Error("Внешний API вернул ошибку")
		return false, err
	}

	var songDetail = models.SongDetail{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Ошибка чтения ответа от внешнего API")
		return false, err
	}

	err = songDetail.UnmarshalJSON(body)
	if err != nil {
		s.log.Error("Ошибка обработки ответа внешнего API")
	}

	success, err := s.songSaver.InsertSong(ctx, song, songDetail)

	return success, nil
}
func (s *Service) GetSong(ctx context.Context,
	songName string,
	groupName string,
	releaseDate string,
	limit int,
	offset int,
) (models.Songs, error) {
	songs, err := s.libraryProvider.Song(ctx, songName, groupName, releaseDate, limit, offset)
	if err != nil {
		return nil, err
	}
	return songs, nil
}
func (s *Service) GetLyric(ctx context.Context,
	songName string,
	groupName string,
	limit int,
	offset int,
) (models.Lyrics, error) {
	lyrics, err := s.libraryProvider.Lyric(ctx, songName, groupName, limit, offset)
	if err != nil {
		return nil, err
	}
	return lyrics, nil
}
func (s *Service) DeleteSong(ctx context.Context,
	songName string,
	groupName string,
) (bool, error) {
	success, err := s.songDeleter.DeleteSong(ctx, songName, groupName)
	if err != nil {
		return success, err
	}
	return success, nil
}
func (s *Service) UpdateSong(ctx context.Context,
	song models.Song,
) (bool, error) {
	success, err := s.songSaver.UpdateSong(ctx, song)
	if err != nil {
		return success, err
	}
	return success, nil
}
