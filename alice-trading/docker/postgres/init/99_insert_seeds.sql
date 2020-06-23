-- 銘柄
insert into instruments (name, pip_location, maximum_order_units, minimum_trade_size, minimum_trailing_stop_distance) values ('USD_JPY', -2, 100000000, 1, 100.000);
insert into instruments (name, pip_location, maximum_order_units, minimum_trade_size, minimum_trailing_stop_distance) values ('EUR_USD', -4, 100000000, 1, 1.00000);
insert into instruments (name, pip_location, maximum_order_units, minimum_trade_size, minimum_trailing_stop_distance) values ('EUR_JPY', -2, 100000000, 1, 0.05);

-- Sequence
insert into sequences (event, sequence) values ('SWING', 1);
