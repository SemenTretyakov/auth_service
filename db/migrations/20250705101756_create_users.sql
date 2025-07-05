-- +goose Up
SELECT 'up SQL query';
CREATE table users (
  id serial primary key,
  fullname text not null,
  createdAt timestamp not null default now(),
  updatedAt timestamp
);


-- +goose Down
SELECT 'down SQL query';
drop table users;
