create table trade_rule_statuses
(
    id serial primary key ,
    trade_rule varchar not null ,
    instrument varchar not null ,
    granularity varchar not null ,
    candle_time timestamptz not null ,
    status bool default true not null
)
