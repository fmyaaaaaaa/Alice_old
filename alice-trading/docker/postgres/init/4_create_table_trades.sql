create table trades
(
    id serial primary key ,
    trade_id int not null unique ,
    units numeric not null ,
    price numeric not null ,
    instrument varchar ,
    state varchar ,
    initial_units numeric ,
    current_units numeric ,
    realized_pl numeric ,
    unrealized_pl numeric ,
    margin_used numeric ,
    open_time timestamptz ,
    close_time timestamptz
)
