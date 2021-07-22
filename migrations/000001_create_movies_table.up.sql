CREATE TABLE IF NOT EXISTS movies
(
    id         bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title      text                        NOT NULL,
    year       integer                     NOT NULL,
    runtime    integer                     NOT NULL,
    genres     text[]                      NOT NULL,
    version    integer                     NOT NULL DEFAULT 1
);


-- bigserial: 64-bit autoincrementing number
-- primary key: Primary Key
-- text[]: Array of Zero or More text values (can be indexed+queried)
-- Not NULL to account for go's idiosyncracies with null values
-- text: Instead of varchar or varchar(n)
-- https://www.depesz.com/2010/03/02/charx-vs-varcharx-vs-varchar-vs-text/