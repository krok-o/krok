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
    url VARCHAR ( 256 )
);

create table rel_command_repositories (
    id serial primary key,
    repository_id int not null,
    command_id int unique not null
);

create table rel_repositories_command (
    id serial primary key,
    command_id int not null,
    repository_id int unique not null
);

create table users (
    id serial primary key,
    -- email is coming from openid registration.
    email varchar(256) unique not null,
    last_login date,
    display_name varchar(50)
);

-- api keys will be generated by the user.
create table apikeys (
    id serial primary key,
    api_key_id text unique not null,
    api_key_secret text not null,
    user_id int not null,
    ttl date
);

-- The files lock which will contain the lock for a file with a timestamp of creation.
-- Locks that are older than 10 minutes will be purged.
create table file_lock (
    name varchar ( 256 ) unique not null,
    lock_start date
)