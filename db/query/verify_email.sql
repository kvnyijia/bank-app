-- name: CreateVerifyEmail :one
insert into 
  verify_emails (
    username,
    email,
    secret_code
  ) 
values (
  $1, $2, $3
) 
returning *
;
