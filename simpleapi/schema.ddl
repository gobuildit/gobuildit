CREATE TABLE developers (
    id serial primary key,
    name varchar,
    age integer
);

INSERT INTO developers (name, age)
VALUES
('Alice', 23),
('Bob', 21);
