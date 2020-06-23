create table swing_targets
(
    id serial primary key ,
    instrument varchar not null ,
    granularity varchar not null ,
    swing_id int not null
)
