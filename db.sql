drop database if exists test_exercise_db;

\connect test_exercise_db

create table users(
d serial primary key,
email varchar not null unique,
name varchar not null,
password varchar not null,
token varchar
);

create table companies(
id serial primary key,
name varchar(15) not null unique,
description varchar(3000),
amount_of_employees int not null,
registered boolean not null,
type varchar not null
);

create table events(
id bigserial primary key,
method varchar not null,
user_email varchar not null,
company_name varchar not null,
created_at timestamp not null default now()
);
