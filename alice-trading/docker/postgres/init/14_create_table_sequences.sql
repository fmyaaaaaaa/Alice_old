create table sequences
(
    id serial primary key ,
    event varchar not null unique ,
    sequence int default 0
)
