create table instruments
(
    id serial primary key ,
    name varchar(50) unique not null ,
    pip_location numeric ,
    maximum_order_units numeric ,
    minimum_trade_size numeric ,
    minimum_trailing_stop_distance numeric
);

insert into instruments (name, pip_location, maximum_order_units, minimum_trade_size, minimum_trailing_stop_distance) values ('USD_JPY', 0.1, 100, 10, 0.05)