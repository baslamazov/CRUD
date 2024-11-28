package services

import (
	"EffectiveMobile/internal/domain/models"
	"context"
	"errors"
	"log/slog"
	"time"
)

type Service struct {
	log          *slog.Logger
	songSaver    SongSaver
	songProvider SongProvider
	tokenTTL     time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type SongSaver interface {
	SaveSong(
		ctx context.Context,
		song string,
		group string,
	) (uid int64, err error)
}

type SongProvider interface {
	Song(ctx context.Context,
		songName string,
		groupName string,
		releaseDate string,
		page int,
		limit int,
		offset int,
	) (models.SongsResponse, error)
}

func New(
	log *slog.Logger,
	songSaver SongSaver,
	songProvider SongProvider,
) *Service {
	return &Service{
		songSaver:    songSaver,
		songProvider: songProvider,
		log:          log,
	}
}
func (s *Service) AddNewSong(ctx context.Context,
	songName string,
	groupName string) (bool, error) {
	return false, nil
}
func (s *Service) GetSong(ctx context.Context,
	songName string,
	groupName string,
	releaseDate string,
	page int,
	limit int,
	offset int,
) (models.SongsResponse, error) {
	songs, err := s.songProvider.Song(ctx, songName, groupName, releaseDate, page, limit, offset)
	if err != nil {
		return nil, err
	}
	return songs, nil
}
