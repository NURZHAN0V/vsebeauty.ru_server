-- Удаляем таблицы в обратном порядке (из-за зависимостей)
DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS mailboxes;