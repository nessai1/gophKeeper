CREATE TABLE IF NOT EXISTS users (
    uuid uuid primary key,
    login varchar(255) not null unique,
    password varchar(255) not null
);