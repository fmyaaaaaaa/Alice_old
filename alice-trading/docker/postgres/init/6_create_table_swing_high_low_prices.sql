create table swing_high_low_prices
(
    id serial primary key ,
    swing_id int not null unique ,
    high_price numeric not null ,
    low_price numeric not null
)
