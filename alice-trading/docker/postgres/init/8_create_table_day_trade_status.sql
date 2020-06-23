create table captain_america
(
    id serial primary key ,
    instrument varchar not null ,
    trend varchar not null ,
    setup_price numeric not null ,
    status bool default false not null
)
