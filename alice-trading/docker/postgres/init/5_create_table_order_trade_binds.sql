create table order_trade_binds
(
    id serial primary key ,
    entry_order_id int not null ,
    trade_id int not null ,
    stop_loss_order_id int not null ,
    is_delete bool default false
)
