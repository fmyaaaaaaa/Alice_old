create table trend_statuses
(
    id serial primary key ,
    instrument varchar not null ,
    granularity varchar not null ,
    trend varchar not null ,
    last_swing_id int not null
)
