-- Подключаем расширение для генерации UUID, если оно еще не включено  
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Вставляем группы
INSERT INTO groups (id, name) VALUES
                                  (gen_random_uuid(), 'Rammstein'),
                                  (gen_random_uuid(), 'AC/DC');

-- Получаем ID групп для использования в дальнейшем
WITH group_ids AS (
    SELECT id, name FROM groups WHERE name IN ('Rammstein', 'AC/DC')
),

-- Вставляем песни
     song_data AS (
         INSERT INTO songs (id, group_id, name, release_date, link)
             SELECT gen_random_uuid(), gid.id, s.name, s.release_date::date, s.link
             FROM (
                      VALUES
                          ('Zeig dich', '2019-05-17', 'https://example.com/rammstein/zeig_dich', 'Rammstein'),
                          ('Reise, Reise', '2004-09-27', 'https://example.com/rammstein/reise_reise', 'Rammstein'),
                          ('Links 2 3 4', '2001-05-14', 'https://example.com/rammstein/links_2_3_4', 'Rammstein'),
                          ('Thunderstruck', '1990-09-10', 'https://example.com/acdc/thunderstruck', 'AC/DC')
                  ) AS s (name, release_date, link, group_name)
                      JOIN group_ids gid ON gid.name = s.group_name
             RETURNING id, group_id, name
     )

-- Подготавливаем данные для вставки текстов песен
SELECT
    sd.id AS song_id,
    sd.group_id AS group_id,
    v.verse_number,
    v.text
INTO TEMP TABLE lyrics_temp
FROM song_data sd
         JOIN (
    -- Куплеты для 'Zeig dich' группы 'Rammstein'
    SELECT 'Zeig dich' AS name, 1 AS verse_number, 'Куплет 1 для Zeig dich' AS text
    UNION ALL
    SELECT 'Zeig dich', 2, 'Куплет 2 для Zeig dich' AS text
    UNION ALL
    SELECT 'Zeig dich', 3, 'Куплет 3 для Zeig dich' AS text
    -- Добавьте остальные куплеты для 'Zeig dich' здесь

    UNION ALL
    -- Куплеты для 'Reise, Reise' группы 'Rammstein'
    SELECT 'Reise, Reise', 1, 'Куплет 1 для Reise, Reise' AS text
    UNION ALL
    SELECT 'Reise, Reise', 2, 'Куплет 2 для Reise, Reise' AS text
    UNION ALL
    SELECT 'Reise, Reise', 3, 'Куплет 3 для Reise, Reise' AS text
    -- Добавьте остальные куплеты для 'Reise, Reise' здесь

    UNION ALL
    -- Куплеты для 'Links 2 3 4' группы 'Rammstein'
    SELECT 'Links 2 3 4', 1, 'Куплет 1 для Links 2 3 4' AS text
    UNION ALL
    SELECT 'Links 2 3 4', 2, 'Куплет 2 для Links 2 3 4' AS text
    UNION ALL
    SELECT 'Links 2 3 4', 3, 'Куплет 3 для Links 2 3 4' AS text
    -- Добавьте остальные куплеты для 'Links 2 3 4' здесь

    UNION ALL
    -- Куплеты для 'Thunderstruck' группы 'AC/DC'
    SELECT 'Thunderstruck', 1, 'Куплет 1 для Thunderstruck' AS text
    UNION ALL
    SELECT 'Thunderstruck', 2, 'Куплет 2 для Thunderstruck' AS text
    UNION ALL
    SELECT 'Thunderstruck', 3, 'Куплет 3 для Thunderstruck' AS text
    -- Добавьте остальные куплеты для 'Thunderstruck' здесь
) v ON sd.name = v.name;

-- Вставляем тексты песен (куплеты)
INSERT INTO lyrics (song_id, group_id, verse_number, text)
SELECT song_id, group_id, verse_number, text FROM lyrics_temp;

-- Удаляем временную таблицу
DROP TABLE lyrics_temp;