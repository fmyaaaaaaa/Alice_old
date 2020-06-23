create table iron_man_statuses
(
    id serial primary key ,
    instrument varchar not null ,
    granularity varchar not null ,
    swing_target_id int not null ,
    trend varchar not null ,
    status bool default true not null
)
