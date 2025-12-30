create table users (
   id         serial primary key,
   username   text not null unique,
   password   text not null,
   role       text not null default 'user',
   avatar     text,
   locked     boolean default false,
   coin       integer default 0,
   created_at timestamptz default now()
);