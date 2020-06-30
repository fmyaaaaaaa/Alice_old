create table accounts
(
    id serial primary key ,
    margin_rate numeric not null ,
    balance numeric not null ,
    open_trade_count numeric ,
    open_position_count numeric ,
    pending_order_count numeric ,
    pl numeric ,
    unrealized_pl numeric not null ,
    nav numeric not null ,
    margin_used numeric not null ,
    margin_available numeric not null ,
    position_value numeric not null ,
    margin_closeout_unrealized_pl numeric not null ,
    margin_closeout_nav numeric not null ,
    margin_closeout_margin_used numeric not null ,
    margin_closeout_position_value numeric not null ,
    margin_closeout_percent numeric not null
)
