-- Создаём таблицу почтовых ящиков
CREATE TABLE IF NOT EXISTS mailboxes (
    id UUID PRIMARY KEY,                           -- Уникальный идентификатор
    address VARCHAR(255) UNIQUE NOT NULL,          -- Email адрес (уникальный)
    created_at TIMESTAMP DEFAULT NOW(),            -- Дата создания
    expires_at TIMESTAMP NOT NULL,                 -- Дата истечения
    is_active BOOLEAN DEFAULT TRUE                 -- Активен ли ящик
);

-- Создаём таблицу писем
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,                           -- Уникальный идентификатор
    mailbox_id UUID REFERENCES mailboxes(id) ON DELETE CASCADE, -- Связь с ящиком
    from_address VARCHAR(255) NOT NULL,            -- Адрес отправителя
    subject VARCHAR(500),                          -- Тема письма
    body_text TEXT,                                -- Текстовое содержимое
    body_html TEXT,                                -- HTML содержимое
    received_at TIMESTAMP DEFAULT NOW(),           -- Дата получения
    is_read BOOLEAN DEFAULT FALSE,                 -- Прочитано ли
    is_spam BOOLEAN DEFAULT FALSE                  -- Спам ли
);

-- Создаём таблицу вложений
CREATE TABLE IF NOT EXISTS attachments (
    id UUID PRIMARY KEY,                           -- Уникальный идентификатор
    message_id UUID REFERENCES messages(id) ON DELETE CASCADE, -- Связь с письмом
    filename VARCHAR(255) NOT NULL,                -- Имя файла
    content_type VARCHAR(100),                     -- MIME-тип
    size_bytes BIGINT,                             -- Размер в байтах
    storage_path VARCHAR(500)                      -- Путь к файлу
);

-- Создаём индексы для ускорения поиска
CREATE INDEX IF NOT EXISTS idx_mailboxes_address ON mailboxes(address);
CREATE INDEX IF NOT EXISTS idx_mailboxes_expires ON mailboxes(expires_at);
CREATE INDEX IF NOT EXISTS idx_messages_mailbox ON messages(mailbox_id);
CREATE INDEX IF NOT EXISTS idx_messages_received ON messages(received_at);