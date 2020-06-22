create table orders
(
    id serial primary key ,
    order_id int not null unique ,
    instrument varchar not null ,
    units numeric not null ,
    type varchar not null ,
    price numeric ,
    distance numeric ,
    time timestamptz not null ,
    commission numeric ,
    time_in_force varchar
)
