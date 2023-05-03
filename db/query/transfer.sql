-- name: CreateTransfer :one
insert into transfers (from_account_id, to_account_id, amount) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTransfer :one
select * 
from transfers
where id = $1
limit 1;