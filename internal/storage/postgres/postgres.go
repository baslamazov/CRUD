package postgres

import (
	"EffectiveMobile/internal/config"
	"EffectiveMobile/internal/domain/models"
	"context"
	"database/sql"
	"errors"
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
	cfg, err := pgxpool.ParseConfig(constr)

	if err != nil {
		log.Fatal(err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	return pool
}

func (db *Storage) Song(ctx context.Context,
	songName string,
	groupName string,
	releaseDate string,
	limit int,
	offset int,
) (models.Songs, error) {
	if err := db.db.Ping(ctx); err != nil {
		return nil, err
	}
	conn, err := db.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`  
        SELECT songs.id, songs.name, songs.release_date, songs.link, groups.name AS group_name  
        FROM songs  
        JOIN groups ON songs.group_id = groups.id  
    `)
	var conditions []string
	var args []interface{}
	argID := 1
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
	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}
	queryBuilder.WriteString(" ORDER BY songs.id")
	if limit <= 0 {
		limit = 10
	}

	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1))
	args = append(args, limit, offset)
	argID += 2

	query := queryBuilder.String()

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs models.Songs

	for rows.Next() {
		var song models.Song
		var releaseDate sql.NullTime
		err := rows.Scan(&song.ID, &song.Name, &releaseDate, &song.Link, &song.GroupName)
		if err != nil {
			return nil, err
		}

		if releaseDate.Valid {
			song.ReleaseDate = releaseDate.Time.Format("2006-01-02")
		} else {
			song.ReleaseDate = ""
		}

		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return songs, nil
}
func (db *Storage) Lyric(ctx context.Context,
	songName string,
	groupName string,
	limit int,
	offset int,
) (models.Lyrics, error) {
	if err := db.db.Ping(ctx); err != nil {
		return nil, err
	}
	conn, err := db.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if songName == "" || groupName == "" {
		return nil, errors.New("songName и groupName должны быть заполнены")
	}

	query := `  
		SELECT l.verse_number, l.text  
		FROM lyrics l  
		JOIN songs s ON l.song_id = s.id  
		JOIN groups g ON s.group_id = g.id  
		WHERE s.name ILIKE $1 AND g.name ILIKE $2  
		ORDER BY l.verse_number  
		LIMIT $3 OFFSET $4  
	`

	rows, err := tx.Query(ctx, query, "%"+songName+"%", "%"+groupName+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lyrics models.Lyrics

	for rows.Next() {
		var lyric models.Lyric
		if err := rows.Scan(&lyric.VerseNumber, &lyric.Text); err != nil {
			return nil, err
		}
		lyrics = append(lyrics, lyric)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return lyrics, nil
}

func (db *Storage) DeleteSong(ctx context.Context,
	songName string,
	groupName string,
) (bool, error) {
	if err := db.db.Ping(ctx); err != nil {
		return false, err
	}
	conn, err := db.db.Acquire(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	query := `  
        DELETE FROM songs  
        USING groups  
        WHERE songs.group_id = groups.id  
          AND songs.name = $1  
          AND groups.name = $2  
    `

	cmdTag, err := tx.Exec(ctx, query, songName, groupName)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	if err = tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (db *Storage) InsertSong(
	ctx context.Context,
	song models.Song,
	songDetail models.SongDetail,
) (bool, error) {
	if err := db.db.Ping(ctx); err != nil {
		return false, err
	}
	conn, err := db.db.Acquire(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	var groupID string
	err = tx.QueryRow(ctx, "SELECT id FROM groups WHERE name = $1", song.GroupName).Scan(&groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = tx.QueryRow(ctx, "INSERT INTO groups (name) VALUES ($1) RETURNING id", song.GroupName).Scan(&groupID)
			if err != nil {
				return false, err
			}
		} else {
			return false, err
		}
	}

	var songID string
	insertSongQuery := `  
        INSERT INTO songs (group_id, name, release_date, link)  
        VALUES ($1, $2, $3, $4)  
        RETURNING id  
    `
	err = tx.QueryRow(ctx, insertSongQuery, groupID, song.Name, songDetail.ReleaseDate, songDetail.Link).Scan(&songID)
	if err != nil {
		return false, err
	}

	verses := strings.Split(songDetail.Text, "\n\n")
	for i, verse := range verses {
		insertLyricQuery := `  
            INSERT INTO lyrics (song_id, group_id, verse_number, text)  
            VALUES ($1, $2, $3, $4)  
        `
		_, err := tx.Exec(ctx, insertLyricQuery, songID, groupID, i+1, verse)
		if err != nil {
			return false, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return false, err
	}
	return true, nil
}
