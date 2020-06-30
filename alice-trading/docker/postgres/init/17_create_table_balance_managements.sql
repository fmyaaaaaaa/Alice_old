create table balance_managements
(
    id serial primary key ,
    trade_id int ,
    instrument varchar ,
    trade_rule varchar ,
    current_account_level numeric ,
    next_account_level numeric ,
    position_size numeric ,
    exec_price numeric ,
    distance numeric ,
    delta numeric ,
    created_at timestamptz
)
