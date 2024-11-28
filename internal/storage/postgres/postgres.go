package postgres

import (
	"EffectiveMobile/internal/config"
	"EffectiveMobile/internal/domain/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewDBService(cfg config.Database) *Storage {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	return &Storage{db: CreatePool(connectionString)}
}

func CreatePool(constr string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(constr)

	if err != nil {
		log.Fatal(err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	return pool
}

func (db *Storage) Song(ctx context.Context,
	songName string,
	groupName string,
	releaseDate string,
	page int,
	limit int,
	offset int,
) (models.SongsResponse, error) {
	// Проверяем соединение с базой данных
	if err := db.db.Ping(ctx); err != nil {
		return nil, err
	}

	// Получаем соединение из пула
	conn, err := db.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	// Начинаем транзакцию
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Строим базовый SQL-запрос для выборки песен
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`  
        SELECT songs.id, songs.name, songs.release_date, songs.link, groups.name AS group_name  
        FROM songs  
        JOIN groups ON songs.group_id = groups.id  
    `)

	// Список условий и аргументов для фильтрации
	var conditions []string
	var args []interface{}
	argID := 1

	// Добавляем фильтры, если параметры не пустые
	if songName != "" {
		conditions = append(conditions, fmt.Sprintf("songs.name ILIKE $%d", argID))
		args = append(args, "%"+songName+"%")
		argID++
	}

	if groupName != "" {
		conditions = append(conditions, fmt.Sprintf("groups.name ILIKE $%d", argID))
		args = append(args, "%"+groupName+"%")
		argID++
	}

	if releaseDate != "" {
		conditions = append(conditions, fmt.Sprintf("songs.release_date = $%d", argID))
		args = append(args, releaseDate)
		argID++
	}

	// Если есть условия, добавляем их в запрос
	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	// Добавляем сортировку по ID песни (при необходимости можно изменить поле сортировки)
	queryBuilder.WriteString(" ORDER BY songs.id")

	// Устанавливаем значения по умолчанию для пагинации
	if limit <= 0 {
		limit = 10 // Значение по умолчанию
	}
	if page > 0 {
		offset = (page - 1) * limit
	}

	// Добавляем параметры LIMIT и OFFSET в запрос
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1))
	args = append(args, limit, offset)
	argID += 2

	// Получаем итоговый запрос
	query := queryBuilder.String()

	// Выполняем запрос для получения песен
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songsResponse models.SongsResponse
	songIDs := []string{}
	songMap := make(map[string]*models.SongResponse)

	// Обрабатываем результаты запроса песен
	for rows.Next() {
		var song models.SongResponse
		var releaseDate sql.NullTime
		err := rows.Scan(&song.ID, &song.Name, &releaseDate, &song.Link, &song.GroupName)
		if err != nil {
			return nil, err
		}

		// Обработка NULL значений даты выпуска
		if releaseDate.Valid {
			song.ReleaseDate = releaseDate.Time.Format("2006-01-02")
		} else {
			song.ReleaseDate = ""
		}

		song.Lyrics = []models.Lyric{}
		songsResponse = append(songsResponse, song)
		songIDs = append(songIDs, song.ID)
		songMap[song.ID] = &songsResponse[len(songsResponse)-1]
	}

	// Проверяем ошибки при итерации результатов
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Если нет песен для обработки, возвращаем пустой результат
	if len(songIDs) == 0 {
		// Фиксируем транзакцию
		if err = tx.Commit(ctx); err != nil {
			return nil, err
		}
		return songsResponse, nil
	}

	// Запрос для получения куплетов песен
	lyricsQuery := `  
        SELECT song_id, verse_number, text  
        FROM lyrics  
        WHERE song_id = ANY(\$1)  
        ORDER BY song_id, verse_number  
    `

	// Выполняем запрос для получения куплетов
	lyricRows, err := tx.Query(ctx, lyricsQuery, songIDs)
	if err != nil {
		return nil, err
	}
	defer lyricRows.Close()

	// Обрабатываем результаты запроса куплетов
	for lyricRows.Next() {
		var songID string
		var lyric models.Lyric
		err := lyricRows.Scan(&songID, &lyric.VerseNumber, &lyric.Text)
		if err != nil {
			return nil, err
		}

		if song, ok := songMap[songID]; ok {
			song.Lyrics = append(song.Lyrics, lyric)
		}
	}

	// Проверяем ошибки при итерации результатов
	if err = lyricRows.Err(); err != nil {
		return nil, err
	}

	// Фиксируем транзакцию
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return songsResponse, nil
}
func (db *Storage) SaveSong(
	ctx context.Context,
	songName string,
	groupName string,
) (int64, error) {
	return int64(1), nil
}
