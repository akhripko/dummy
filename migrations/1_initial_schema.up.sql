BEGIN;

create table activity_log (
    name varchar(36) not null,
    last_touch timestamp not null,
    constraint activity_log_pk
        primary key (name)
);

COMMIT;