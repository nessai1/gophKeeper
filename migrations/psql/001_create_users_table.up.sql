CREATE TABLE IF NOT EXISTS users (
    id serial primary key,
    login varchar(255) not null unique,
    password varchar(255) not null
);