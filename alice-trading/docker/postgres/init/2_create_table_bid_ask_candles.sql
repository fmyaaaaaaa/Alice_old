create table bid_ask_candles
(
    id              serial primary key,
    instrument_name varchar,
    granularity     varchar,
    open_bid        numeric,
    open_ask        numeric,
    close_bid       numeric,
    close_ask       numeric,
    high_bid        numeric,
    high_ask        numeric,
    low_bid         numeric,
    low_ask         numeric,
    time            timestamptz,
    volume          numeric,
    line            varchar,
    trend           varchar,
    swing_id        int
);
