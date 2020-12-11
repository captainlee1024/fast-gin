CREATE TABLE fast_gin(
    id              bigint auto_increment               primary key COMMENT '自增id',
    fast_gin_id     bigint                              not null COMMENT '用户id',
    demo_name       varchar(64)                         not null COMMENT '姓名',
    info            varchar(64)                         not null COMMENT '详细信息',
    create_time     timestamp default CURRENT_TIMESTAMP null COMMENT '新增时间',
    update_time     timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP COMMENT '更新时间',
    constraint idx_fast_gin_id
        unique (fast_gin_id),
    constraint idx_username
        unique (demo_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='fast_gin测试表';