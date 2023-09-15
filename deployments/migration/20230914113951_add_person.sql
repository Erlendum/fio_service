-- +goose Up
-- +goose StatementBegin
create schema if not exists service;
create table service.persons (
    id serial primary key,
    name text not null,
    surname text not null,
    patronymic text,
    age int not null,
    gender text not null,
    nationality text not null
);

alter table service.persons
    add constraint correct_gender check ( gender = 'Male' or gender = 'Female' );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists service.persons;
drop schema if exists service;
-- +goose StatementEnd
