create table back_test_results
(
    id serial primary key ,
    trade_rule varchar not null ,
    instrument varchar not null ,
    win_rate numeric not null ,
    lose_rate numeric not null ,
    max_draw_down numeric not null
)
