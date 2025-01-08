CREATE TABLE `nav_user_account` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `user_name` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名称',
    `password` varchar(50) NOT NULL DEFAULT '' COMMENT '密码',
    `account_type` varchar(15) NOT NULL DEFAULT 'private' COMMENT '账号类型,public or private',
    `create_user` varchar(50) NOT NULL DEFAULT '' COMMENT '创建者',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `user_name` (`user_name`),
    UNIQUE (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8 COMMENT='用户账号列表';
CREATE TABLE `nav_user` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `user_name` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名称',
    `user_id` varchar(50) NOT NULL DEFAULT '' COMMENT '用户唯一id',
    `user_type` varchar(2) NOT NULL DEFAULT '' COMMENT '用户类型',
    `user_company` varchar(50) NOT NULL DEFAULT '' COMMENT '用户公司名称',
    `actions` varchar(200) NOT NULL DEFAULT '' COMMENT '服务类别,权限点',
    `ship_name` varchar(100) NOT NULL DEFAULT '' COMMENT '船舶名称',
    `ship_id` varchar(50) NOT NULL DEFAULT '' COMMENT '船舶id,该用户可看哪艘船',
    `ship_group` varchar(100) NOT NULL DEFAULT '' COMMENT '船舶分组',
    `begin_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '起始时间',
    `end_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '结束时间',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '更新时间',
    -- `sign` varchar(200) NOT NULL DEFAULT '' COMMENT '盐值',
    PRIMARY KEY (`id`),
    KEY `user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8 COMMENT='用户列表';
CREATE TABLE `nav_ships` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `ship_group` varchar(100) NOT NULL DEFAULT '' COMMENT '船舶分组',
  `ship_manage_company` varchar(100) NOT NULL DEFAULT '瑞宁航运' COMMENT '船舶管理公司',
  `ship_manufacturer` varchar(100) NOT NULL DEFAULT '' COMMENT '船舶生产厂家',
  `ship_name` varchar(100) NOT NULL DEFAULT '' COMMENT '船舶名称',
  `ship_type` varchar(100) NOT NULL DEFAULT '' COMMENT '船舶型号',
  `ship_id` varchar(100) NOT NULL DEFAULT '' COMMENT '船舶唯一标识',
  `ship_img` varchar(300) NOT NULL DEFAULT '' COMMENT '船舶图片',
  `captain_name` varchar(50) NOT NULL DEFAULT '' COMMENT '机务主管姓名',
  `captain_phone` varchar(50) NOT NULL DEFAULT '' COMMENT '机务手机号',
  `captain_mail` varchar(50) NOT NULL DEFAULT '' COMMENT '机务邮箱',
  `from` varchar(300) NOT NULL DEFAULT '' COMMENT '出发地点',
  `destination` varchar(300) NOT NULL DEFAULT '' COMMENT '目的地',
  `leave_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '离港时间',
  `task_desc` varchar(300) NOT NULL DEFAULT '粮油运输' COMMENT '作业描述',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '更新时间',
  `is_remove` varchar(3) NOT NULL DEFAULT '0' COMMENT '船舶是否移除',
  PRIMARY KEY (`id`),
  KEY `ship_id` (`ship_id`, `is_remove`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8 COMMENT='船舶列表';

CREATE TABLE `nav_status_gongyin1` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `ship_id` varchar(100) NOT NULL DEFAULT '' COMMENT '船舶唯一标识',
  `channel_id` varchar(50) NOT NULL DEFAULT '' COMMENT '监测点id',
  `channel_attribute` varchar(50) NOT NULL DEFAULT 'DI' COMMENT '通道属性AI/DI',
  `channel_description` varchar(50) NOT NULL DEFAULT '' COMMENT '监测点描述',
  `channel_status` varchar(50) NOT NULL DEFAULT '' COMMENT '监测点状态',
  `channel_value` varchar(50) NOT NULL DEFAULT '' COMMENT '监测点数值',
  `channel_unit` varchar(10) NOT NULL DEFAULT '' COMMENT '单位',
  `user_alarm_value` int(10) NOT NULL DEFAULT '0' COMMENT '报警阈值,岸上设置',
  `user_alarm_status` varchar(10) NOT NULL DEFAULT 'Normal' COMMENT '报警状态,基于岸上设置报警值',
  `channel_time` varchar(25) NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '数据产生时间',
  `file_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '文件生成时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `is_remove` varchar(3) NOT NULL DEFAULT '0' COMMENT '通道是否移除展示',
  PRIMARY KEY (`id`),
  KEY `ship_id` (`ship_id`, `create_time`, `channel_id`, `channel_status`, `user_alarm_status`,`channel_attribute`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8 COMMENT='船舶数据表';

CREATE TABLE `nav_status_ruining2` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `ship_id` varchar(100) NOT NULL DEFAULT '' COMMENT '船舶唯一标识',
  `channel_id` varchar(50) NOT NULL DEFAULT '' COMMENT '监测点id',
  `channel_attribute` varchar(50) NOT NULL DEFAULT 'DI' COMMENT '通道属性AI/DI',
  `channel_description` varchar(50) NOT NULL DEFAULT '' COMMENT '监测点描述',
  `channel_status` varchar(50) NOT NULL DEFAULT '' COMMENT '监测点状态',
  `channel_value` varchar(50) NOT NULL DEFAULT '' COMMENT '监测点数值',
  `channel_unit` varchar(10) NOT NULL DEFAULT '' COMMENT '单位',
  `user_alarm_value` int(10) NOT NULL DEFAULT '0' COMMENT '报警阈值,岸上设置',
  `user_alarm_status` varchar(10) NOT NULL DEFAULT 'Normal' COMMENT '报警状态,基于岸上设置报警值',
  `channel_time` varchar(25) NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '数据产生时间',
  `file_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '文件生成时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `is_remove` varchar(3) NOT NULL DEFAULT '0' COMMENT '通道是否移除展示',
  PRIMARY KEY (`id`),
  KEY `ship_id` (`ship_id`, `create_time`, `channel_id`, `channel_status`, `user_alarm_status`,`channel_attribute`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8 COMMENT='船舶数据表';