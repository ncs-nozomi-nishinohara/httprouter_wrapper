create table if not exists usermaster
(
     id     serial
    ,name   varchar(20)
);
comment on table usermaster is 'ユーザー';
comment on column usermaster.id is 'id';
comment on column usermaster.name is '姓名';
insert into usermaster(name) values ('test');