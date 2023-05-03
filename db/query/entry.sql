-- name: CreateEntry :one
insert into entries (account_id, amount)
values ($1, $2)
RETURNING *;

-- name: GetEntry :one
select * 
from entries
where id = $1
limit 1;