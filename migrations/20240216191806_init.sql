-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS talks (
    url text PRIMARY KEY NOT NULL,
    name text NOT NULL,
    duration interval NOT NULL,
    starts_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS talk_speakers (
    talk_url text NOT NULL REFERENCES talks (url),
    speaker varchar(128) NOT NULL,
    company varchar(255)
);

CREATE TABLE IF NOT EXISTS users (
    telegram_id bigint PRIMARY KEY NOT NULL,
    identification varchar(255),
    favorite_talks text[]
);

CREATE TYPE feedback_type AS ENUM ('Оценка', 'Я не слушал этот доклад', 'Не хочу оценивать');

CREATE TABLE IF NOT EXISTS talk_feedbacks (
    talk_url text NOT NULL,
    telegram_id bigint NOT NULL,
    type feedback_type NOT NULL,
    content_rating smallint,
    performance_rating smallint,
    comment text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS talk_feedbacks;
DROP TYPE IF EXISTS feedback_type;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS talk_speakers;
DROP TABLE IF EXISTS talks;
-- +goose StatementEnd
