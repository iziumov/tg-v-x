CREATE TYPE job_status as ENUM('pending', 'downloading', 'done', 'failed');

CREATE TABLE jobs(
  id serial primary key,
  user_id bigint references users(tg_id) on delete cascade,
  url text not null,
  status job_status default 'pending',
  platform text,
  file_size_bytes bigint,
  error_msg text,
  created_at timestamptz default current_timestamp,
  finished_at timestamptz
);