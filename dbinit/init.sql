create table commands (
    id serial primary key ,
    name varchar unique not null,
    schedule varchar,
    filename varchar unique not null,
    hash varchar unique not null,
    location varchar not null,
    enabled boolean not null
);

create table command_settings
(
    id serial primary key,
    command_id int,
    constraint fk_command_id
        foreign key (command_id)
            references commands(id)
            on delete cascade,
    -- this will have to be appended with the command ID and a unique id
    -- in case it's in_vault to not clash with other settings.
    key varchar,
    value varchar,
    in_vault boolean,
    -- for a command make sure a key is unique. But for other commands the same key can be used.
    unique(command_id, key)
);

create table repositories (
    id serial primary key,
    name varchar ( 256 ) unique not null,
    url varchar ( 256 ),
    vcs int,
    project_id int null
);

create table rel_commands_repositories (
    id serial primary key,
    repository_id int,
    command_id int,
    constraint fk_repository_id
        foreign key (repository_id)
            references repositories(id)
            on delete cascade,
    constraint fk_command_id
        foreign key (command_id)
            references commands(id)
            on delete cascade
);

-- The relationship which defines if a command supports a given platform or not.
-- platform_id is a hardcoded value and only defined in Krok.
-- It won't be something that is configurable. More will be added as more
-- platforms start to be supported.
create table rel_commands_platforms (
    id serial primary key,
    platform_id int,
    command_id int,
    constraint fk_command_id
        foreign key (command_id)
            references commands(id)
            on delete cascade
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
    name varchar,
    api_key_id varchar unique not null,
    -- this will be shown once then never again as it will be stored encrypted.
    api_key_secret varchar not null,
    user_id int not null,
    ttl varchar,
    created_at date
);

-- The files lock which will contain the lock for a file with a timestamp of creation.
-- Locks that are older than 10 minutes will be purged.
-- Note: Delete this when we remove the watcher.
create table file_lock (
    name varchar ( 256 ) unique not null,
    lock_start date
);

-- store events for a repository.
create table events (
    id serial primary key,
    event_id varchar unique not null,
    repository_id int,
    payload varchar,
    created_at date,
    vcs int
);

-- store a run for a command. This is associated with an event.
-- Note that we don't save command ID here, because it might have
-- been already deleted when we look back to this event.
-- So events will contain CommandRuns which save the name of the
-- command only.
create table command_run (
    id serial primary key,
    command_name varchar,
    event_id int,
    status varchar,
    outcome varchar,
    created_at date
);

-- generate default admin user
insert into users (email, last_login, display_name) values ('admin@admin.com', now(), 'Admin');
-- secret is 'secret'
insert into apikeys (name, api_key_id, api_key_secret, user_id, ttl, created_at) values ('test', 'api-key-id', '$2y$12$qu2jd67X2dWJJZHccKPY1O/SB1pQQ/HNpYQiSUGBKjzYWIomZeVmG', 1, '3120h', now());
