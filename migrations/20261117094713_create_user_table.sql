-- +goose Up
CREATE TABLE users (
    id bigserial primary key,
    username varchar(50) not null,
    password_hash varchar(32) not null,
    created_at timestamp(0) with time zone not null
);

-- +goose Down
DROP TABLE users;
