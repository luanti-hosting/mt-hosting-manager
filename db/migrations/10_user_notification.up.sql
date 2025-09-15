
create table user_notification(
    id varchar(36) primary key not null, -- uuid
    user_id varchar(36) not null references public.user(id) on delete cascade,
    timestamp bigint not null, -- in `time.Now().Unix()`
    seen boolean not null default false,
    title varchar not null,
    message varchar not null
);
