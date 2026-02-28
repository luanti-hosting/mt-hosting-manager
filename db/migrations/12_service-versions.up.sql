
alter table minetest_server alter column ui_version type varchar(256);
alter table minetest_server alter column ui_version set default 'v1.87';

alter table minetest_server add column nginx_version varchar(256) not null default '1.25.2';

alter table user_node add column node_exporter_version varchar(256) not null default 'v1.6.1';
alter table user_node add column traefik_version varchar(256) not null default 'v3.1';
alter table user_node add column ipv6nat_version varchar(256) not null default '0.4.4';

create table image_version(
    name varchar(256) primary key not null, -- nginx, mtui
    version varchar(256) not null -- 2.1.2
);

insert into image_version(name, version) values('mtui', 'v1.87');
insert into image_version(name, version) values('nginx', '1.25.2');
insert into image_version(name, version) values('node_exporter', 'v1.6.1');
insert into image_version(name, version) values('traefik', 'v3.1');
insert into image_version(name, version) values('ipv6nat', '0.4.4');