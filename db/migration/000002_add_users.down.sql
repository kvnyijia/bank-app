drop index if exists "accounts_owner_currency_idx";
alter table if exists "accounts" drop constraint if exists "accounts_owner_fkey";
drop table if exists "users";