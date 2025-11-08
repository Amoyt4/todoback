-- init.sql
-- Инициализация базы данных для Telegram бота заметок

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     telegram_id BIGINT UNIQUE NOT NULL,
                                     username TEXT,
                                     first_name TEXT,
                                     last_name TEXT,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица категорий
CREATE TABLE IF NOT EXISTS categories (
                                          id SERIAL PRIMARY KEY,
                                          user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                          name TEXT NOT NULL,
                                          color TEXT DEFAULT '#007ACC',
                                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                          UNIQUE(user_id, name)
);

-- Таблица заметок
CREATE TABLE IF NOT EXISTS notes (
                                     id SERIAL PRIMARY KEY,
                                     user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                     title TEXT NOT NULL DEFAULT 'Без названия',
                                     content TEXT,
                                     category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
                                     is_pinned BOOLEAN DEFAULT FALSE,
                                     is_archived BOOLEAN DEFAULT FALSE,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица тегов
CREATE TABLE IF NOT EXISTS tags (
                                    id SERIAL PRIMARY KEY,
                                    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                    name TEXT NOT NULL,
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    UNIQUE(user_id, name)
);

-- Связь многие-ко-многим между заметками и тегами
CREATE TABLE IF NOT EXISTS note_tags (
                                         note_id INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
                                         tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
                                         PRIMARY KEY (note_id, tag_id)
);