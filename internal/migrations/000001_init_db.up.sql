create table if not exists users (
    id bigserial,
    user_name varchar(30) not null,
    first_name varchar(30) not null,
    last_name varchar(30) not null,
    password_hash varchar(72) not null,
    status smallint default 1,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp,
    constraint pk__users primary key(id),
    constraint chk__users__status check(status in (0, 1))
);

create unique index idx__users__user_name on users(user_name) where (deleted_at is null);

create table if not exists secrets (
    id bigserial,
    user_id bigint not null,
    title varchar(100) not null,
    data jsonb not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    constraint pk__secrets primary key(id),
    constraint fk__secrets__user foreign key(user_id) references users(id) on delete cascade
);