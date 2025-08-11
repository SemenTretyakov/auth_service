-- +goose Up
SELECT 'up SQL query';
CREATE table users (
  id serial primary key,
  fullname text not null,
  email text not null,
  password text not null,
  password_confirm text not null,
  role int not null,
  created_at timestamp not null default now(),
  updated_at timestamp
);


-- +goose Down
SELECT 'down SQL query';
drop table users;
