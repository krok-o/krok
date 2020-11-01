CREATE TABLE commands (
    id serial PRIMARY KEY,
    name VARCHAR ( 50 ) UNIQUE NOT NULL,
    schedule VARCHAR ( 50 ),
    filename VARCHAR ( 50 ) UNIQUE NOT NULL,
    hash VARCHAR ( 50 ) UNIQUE NOT NULL,
    location VARCHAR ( 50 ) UNIQUE NOT NULL,
    enabled BOOLEAN NOT NULL
);

CREATE TABLE repositories (
    id serial PRIMARY KEY,
    name VARCHAR ( 256 ) UNIQUE NOT NULL,
    url VARCHAR ( 256 ),
    ssh TEXT NULL
);

create table rel_command_repositories (
    id serial primary key,
    command_id int unique not null,
    repository_id int not null
);

create table rel_repositories_command (
    id serial primary key,
    repository_id int unique not null,
    command_id int not null
);

create table users (
    id serial primary key,
    username varchar(50) unique not null
);