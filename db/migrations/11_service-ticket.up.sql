
create table service_ticket(
    id varchar(36) primary key not null, -- uuid
    title varchar(512) not null,
    user_id varchar(36) not null references public.user(id) on delete cascade,
    user_node_id varchar(36) references public.user_node(id) on delete cascade,
    minetest_server_id varchar(36) references public.minetest_server(id) on delete cascade,
    backup_id varchar(36) references public.backup(id) on delete cascade,
    created bigint not null, -- in `time.Now().Unix()`
    closed bigint, -- in `time.Now().Unix()`
    state varchar not null -- OPEN,RESOLVED,INVALID
);

create table service_ticket_message(
    id varchar(36) primary key not null, -- uuid
    ticket_id varchar(36) not null references public.service_ticket(id) on delete cascade,
    timestamp bigint not null, -- in `time.Now().Unix()`
    message varchar not null
);
