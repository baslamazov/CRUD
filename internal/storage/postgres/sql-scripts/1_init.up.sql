-- Создаем базу данных users  
CREATE DATABASE EffectiveMobile;

\c EffectiveMobile;

-- Удаляем таблицы, если они существуют
DROP TABLE IF EXISTS songs;
DROP TABLE IF EXISTS "groups";
DROP TABLE IF EXISTS lyrics;

-- Включаем расширение pgcrypto для генерации UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Создаем таблицу groups
CREATE TABLE groups (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        name TEXT NOT NULL UNIQUE  -- Название группы должно быть уникальным
);

-- Создаем таблицу songs
CREATE TABLE songs (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
                       name TEXT NOT NULL,
                       release_date DATE,
                       link TEXT
);

-- Создаем таблицу lyrics
CREATE TABLE lyrics (
                        song_id UUID NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
                        group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
                        verse_number INTEGER NOT NULL,
                        text TEXT NOT NULL,
                        PRIMARY KEY (song_id, group_id, verse_number)  -- Составной первичный ключ
);