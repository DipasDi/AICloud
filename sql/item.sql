-- Расширение pgvector (один раз на базу)
CREATE EXTENSION IF NOT EXISTS vector;
CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    target_type TEXT NOT NULL CHECK (
        target_type IN ('file', 'message', 'photo', 'link')
    ),
    tg_file_id TEXT NOT NULL DEFAULT '',
    tg_message_id BIGINT NOT NULL DEFAULT 0,
    file_name TEXT NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    group_name TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Векторный поиск
    embedding vector(4096),
    -- NULL, пока вектор не посчитан
    embedding_model TEXT,
    -- какой моделью посчитан
    -- Полнотекстовый поиск (считается автоматически)
    search_vector TSVECTOR GENERATED ALWAYS AS (
        setweight(
            to_tsvector('russian', coalesce(comment, '')),
            'A'
        ) || setweight(
            to_tsvector(
                'russian',
                coalesce(file_name, '') || ' ' || coalesce(group_name, '')
            ),
            'B'
        )
    ) STORED
);
-- Фильтрация
CREATE INDEX IF NOT EXISTS idx_items_user_chat ON items(user_id, chat_id);
CREATE INDEX IF NOT EXISTS idx_items_group ON items(user_id, group_name);
CREATE INDEX IF NOT EXISTS idx_items_type ON items(user_id, target_type);
CREATE UNIQUE INDEX IF NOT EXISTS uq_items_chat_msg ON items(chat_id, tg_message_id)
WHERE tg_message_id <> 0;
-- Поиск
CREATE INDEX IF NOT EXISTS idx_items_search ON items USING GIN(search_vector);
CREATE INDEX IF NOT EXISTS idx_items_embedding ON items USING hnsw (embedding vector_cosine_ops);
-- Автообновление updated_at
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$ BEGIN NEW.updated_at = NOW();
RETURN NEW;
END $$ LANGUAGE plpgsql;
DROP TRIGGER IF EXISTS trg_items_updated_at ON items;
CREATE TRIGGER trg_items_updated_at BEFORE
UPDATE ON items FOR EACH ROW EXECUTE FUNCTION set_updated_at();