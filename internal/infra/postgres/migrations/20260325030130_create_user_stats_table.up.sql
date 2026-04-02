CREATE TABLE user_stats(
  user_id bigint primary key references users(tg_id) on delete cascade,
  total_jobs int default 0,
  success_jobs int default 0,
  failed_jobs int default 0,
  total_bytes bigint default 0,
  last_active timestamptz default current_timestamp
);