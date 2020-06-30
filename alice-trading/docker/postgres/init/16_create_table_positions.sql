create table positions
(
    id serial primary key ,
    instrument varchar unique ,
    pl numeric not null ,
    unrealized_pl numeric not null ,
    margin_used numeric not null ,
    units numeric not null ,
    tradable_short bool default false ,
    tradable_long bool default false
)
