create table instruments
(
    id serial primary key ,
    name varchar(50) unique not null ,
    pip_location numeric ,
    maximum_order_units numeric ,
    minimum_trade_size numeric ,
    minimum_trailing_stop_distance numeric
);
