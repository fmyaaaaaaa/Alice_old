create table captain_america_statuses
(
    id serial primary key ,
    instrument varchar not null ,
    granularity varchar not null ,
    line varchar not null ,
    setup_price numeric not null ,
    setup_status bool default false not null ,
    trade_status bool default false not null ,
    second_judge bool default false not null
)
