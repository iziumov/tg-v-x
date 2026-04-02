CREATE TABLE users(
  tg_id bigint primary key unique not null,
  username text,
  first_name text,
  last_name text,
  created_at timestamptz default current_timestamp,
  is_admin boolean default false,
  is_banned boolean default false
);