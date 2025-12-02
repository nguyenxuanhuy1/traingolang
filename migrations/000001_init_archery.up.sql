create table users (
   id         serial primary key,
   username   text unique not null,
   avatar     text,
   elo        int default 1000,
   created_at timestamp default now()
);

create table matches (
   id          serial primary key,
   status      text not null,
   max_players int not null,
   created_at  timestamp default now(),
   started_at  timestamp,
   ended_at    timestamp
);

create table match_players (
   id         serial primary key,
   match_id   int
      references matches ( id )
         on delete cascade,
   user_id    int
      references users ( id )
         on delete cascade,
   hp         int not null default 100,
   alive      boolean default true,
   kills      int default 0,
   deaths     int default 0,
   damage     int default 0,
   final_rank int,
   is_winner  boolean default false,
   joined_at  timestamp default now(),
   left_at    timestamp
);