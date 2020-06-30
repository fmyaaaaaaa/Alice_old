create table instruments
(
    id serial primary key ,
    instrument varchar unique not null ,
    pip_location numeric ,
    minimum_trade_size numeric ,
    maximum_trailing_stop_distance numeric ,
    minimum_trailing_stop_distance numeric ,
    maximum_position_size numeric ,
    maximum_order_units numeric ,
    margin_rate numeric ,
    evaluation_instrument varchar ,
    tradable bool
);
