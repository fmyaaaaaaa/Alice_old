-- 銘柄
insert into instruments (instrument, pip_location, minimum_trade_size, maximum_trailing_stop_distance, minimum_trailing_stop_distance, maximum_position_size, maximum_order_units, margin_rate, evaluation_instrument, tradable) values ('USD_JPY', -2, 1, 100.000, 0.050, 0, 100000000, 0.04, 'JPY', true);
insert into instruments (instrument, pip_location, minimum_trade_size, maximum_trailing_stop_distance, minimum_trailing_stop_distance, maximum_position_size, maximum_order_units, margin_rate, evaluation_instrument, tradable) values ('EUR_USD', -4, 1, 1.00000, 0.00050, 0, 100000000, 0.04, 'USD_JPY', true);
insert into instruments (instrument, pip_location, minimum_trade_size, maximum_trailing_stop_distance, minimum_trailing_stop_distance, maximum_position_size, maximum_order_units, margin_rate, evaluation_instrument, tradable) values ('EUR_JPY', -2, 1, 100.000, 0.050, 0, 100000000, 0.04, 'JPY', true);
insert into instruments (instrument, pip_location, minimum_trade_size, maximum_trailing_stop_distance, minimum_trailing_stop_distance, maximum_position_size, maximum_order_units, margin_rate, evaluation_instrument, tradable) values ('AUD_USD', -4, 1, 1.00000, 0.00050, 0, 100000000, 0.04, 'USD_JPY', true);
insert into instruments (instrument, pip_location, minimum_trade_size, maximum_trailing_stop_distance, minimum_trailing_stop_distance, maximum_position_size, maximum_order_units, margin_rate, evaluation_instrument, tradable) values ('EUR_CHF', -4, 1, 1.00000, 0.00050, 0, 100000000, 0.04, 'CHF_JPY', true);
insert into instruments (instrument, pip_location, minimum_trade_size, maximum_trailing_stop_distance, minimum_trailing_stop_distance, maximum_position_size, maximum_order_units, margin_rate, evaluation_instrument, tradable) values ('CHF_JPY', -2, 1, 100.000, 0.050, 0, 100000000, 0.05, 'JPY', true);
insert into instruments (instrument, pip_location, minimum_trade_size, maximum_trailing_stop_distance, minimum_trailing_stop_distance, maximum_position_size, maximum_order_units, margin_rate, evaluation_instrument, tradable) values ('AUD_JPY', -2, 1, 100.000, 0.050, 0, 100000000, 0.04, 'JPY', true);
insert into instruments (instrument, pip_location, minimum_trade_size, maximum_trailing_stop_distance, minimum_trailing_stop_distance, maximum_position_size, maximum_order_units, margin_rate, evaluation_instrument, tradable) values ('GBP_JPY', -2, 1, 100.000, 0.050, 0, 500000000, 0.05, 'JPY', true);
insert into instruments (instrument, pip_location, minimum_trade_size, maximum_trailing_stop_distance, minimum_trailing_stop_distance, maximum_position_size, maximum_order_units, margin_rate, evaluation_instrument, tradable) values ('USD_CAD', -4, 1, 1.00000, 0.00050, 0, 100000000, 0.04, 'CAD_JPY', true);
insert into instruments (instrument, pip_location, minimum_trade_size, maximum_trailing_stop_distance, minimum_trailing_stop_distance, maximum_position_size, maximum_order_units, margin_rate, evaluation_instrument, tradable) values ('NZD_USD', -4, 1, 1.00000, 0.00050, 0, 100000000, 0.04, 'USD_JPY', true);

-- アカウント
insert into accounts (margin_rate, balance, open_trade_count, open_position_count, pending_order_count, pl, unrealized_pl, nav, margin_used, margin_available, position_value, margin_closeout_unrealized_pl, margin_closeout_nav, margin_closeout_margin_used, margin_closeout_position_value, margin_closeout_percent) values (0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0);

-- バックテストリザルト
insert into back_test_results (trade_rule, instrument, win_rate, lose_rate, max_draw_down) values ('CAPTAIN_AMERICA', 'USD_JPY', 0.06, 0.04, 1500);
insert into back_test_results (trade_rule, instrument, win_rate, lose_rate, max_draw_down) values ('CAPTAIN_AMERICA', 'EUR_USD', 0.06, 0.04, 1500);
insert into back_test_results (trade_rule, instrument, win_rate, lose_rate, max_draw_down) values ('CAPTAIN_AMERICA', 'EUR_JPY', 0.06, 0.04, 1500);

-- Sequence
insert into sequences (event, sequence) values ('SWING', 1);