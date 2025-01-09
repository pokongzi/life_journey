
CREATE TABLE IF NOT EXISTS `user` ( 
    `id` BIGINT NOT NULL auto_increment  PRIMARY KEY, 
    `uid` CHAR(20)  NOT NULL DEFAULT '用户id' COMMENT 'uid',
    `name` VARCHAR ( 30 ) NOT NULL DEFAULT '' COMMENT '姓名',
    `birthday` VARCHAR ( 30 ) NOT NULL DEFAULT '' COMMENT '生日',
    `phone` VARCHAR ( 30 ) NOT NULL DEFAULT '' COMMENT '手机号',
    `birthday_type` int NOT NULL DEFAULT 0 COMMENT '类型1阳历2阴历',
    `icon` VARCHAR ( 100 ) NOT NULL DEFAULT '' COMMENT '头像地址',
    `sex` int NOT NULL DEFAULT 0 COMMENT '性别 1男2女',
    `status` int NOT NULL DEFAULT 0 COMMENT '状态 0正常1注销',
    `create_time` int64 NOT NULL DEFAULT 0 COMMENT '创建日期',
    `update_time` int64 NOT NULL DEFAULT 0 COMMENT '更新日期',
) COMMENT '用户信息'


CREATE TABLE IF NOT EXISTS `article` ( 
    `id` BIGINT NOT NULL auto_increment  PRIMARY KEY, 
    `content` VARCHAR NOT NULL DEFAULT '' COMMENT '内容',
    `pic` VARCHAR NOT NULL DEFAULT '' COMMENT '图片',
    `tag` VARCHAR NOT NULL DEFAULT '' COMMENT '标签id',
    `status` int NOT NULL DEFAULT 0 COMMENT '状态 0正常删除',
    `create_time` int64 NOT NULL DEFAULT 0 COMMENT '创建日期',
    `update_time` int64 NOT NULL DEFAULT 0 COMMENT '更新日期',
) COMMENT '文章表'

CREATE TABLE IF NOT EXISTS `relation` ( 
    `id` BIGINT NOT NULL auto_increment  PRIMARY KEY, 
    `uid` CHAR(20)  NOT NULL DEFAULT '用户id' COMMENT 'uid',
    `hashid` CHAR(20)  NOT NULL DEFAULT '哈希id',
    `type` int NOT NULL DEFAULT 0 COMMENT '类型 1 添加者 2关注者',
    `status` int NOT NULL DEFAULT 0 COMMENT '状态 0正常删除',
    `create_time` int64 NOT NULL DEFAULT 0 COMMENT '创建日期',
    `update_time` int64 NOT NULL DEFAULT 0 COMMENT '更新日期',
) COMMENT '关系表'



CREATE TABLE IF NOT EXISTS `like` ( 
    `id` BIGINT NOT NULL auto_increment  PRIMARY KEY, 
    `article_id` int  NOT NULL DEFAULT 0 COMMENT '文章id',
    `uid` CHAR(20)  NOT NULL DEFAULT '用户id' COMMENT 'uid',
    `type` int NOT NULL DEFAULT 0 COMMENT '类型 1点赞 2评论',
    `content` VARCHAR NOT NULL DEFAULT '' COMMENT '内容',
    `status` int NOT NULL DEFAULT 0 COMMENT '状态 0正常删除',
    `create_time` int64 NOT NULL DEFAULT 0 COMMENT '创建日期',
) COMMENT '评论表'


CREATE TABLE IF NOT EXISTS `tag` ( 
    `id` BIGINT NOT NULL auto_increment  PRIMARY KEY, 
    `uid` CHAR(20)  NOT NULL DEFAULT '用户id' COMMENT 'uid',
    `type` int NOT NULL DEFAULT 0 COMMENT '类型 1点赞 2评论',
    `content` VARCHAR NOT NULL DEFAULT '' COMMENT '内容',
    `status` int NOT NULL DEFAULT 0 COMMENT '状态 0正常删除',
    `create_time` int64 NOT NULL DEFAULT 0 COMMENT '创建日期',
) COMMENT '标签表'


CREATE TABLE IF NOT EXISTS `import_day` ( 
    `id` BIGINT NOT NULL auto_increment  PRIMARY KEY, 
    `uid` CHAR(20)  NOT NULL DEFAULT '用户id' COMMENT 'uid',
    `name` VARCHAR ( 30 ) NOT NULL DEFAULT '' COMMENT '姓名',
    `birthday` VARCHAR ( 30 ) NOT NULL DEFAULT '' COMMENT '生日',
    `type` VARCHAR ( 30 ) NOT NULL DEFAULT '' COMMENT '类型 1生日 ',
    `birthday_type` int NOT NULL DEFAULT 0 COMMENT '类型1阳历2阴历',
    `status` int NOT NULL DEFAULT 0 COMMENT '状态 0正常删除',
    `create_time` int64 NOT NULL DEFAULT 0 COMMENT '创建日期',
) COMMENT '表'


CREATE TABLE IF NOT EXISTS `version` ( 
    `id` BIGINT NOT NULL auto_increment  PRIMARY KEY, 
    `name` VARCHAR ( 30 ) NOT NULL DEFAULT '' COMMENT '名称',
    `version` VARCHAR ( 30 ) NOT NULL DEFAULT '' COMMENT '版本',
    `level` int NOT NULL DEFAULT 0 COMMENT '等级',
    `parent_level` int NOT NULL DEFAULT 0 COMMENT '父层级',
    `url` int64 NOT NULL DEFAULT 0 COMMENT '文件地址',
    `tag` text NOT NULL DEFAULT '' COMMENT '版本说明',
    `status` int NOT NULL DEFAULT 0 COMMENT '状态 0正常删除',
    `create_time` int64 NOT NULL DEFAULT 0 COMMENT '创建日期',
    `update_time` int64 NOT NULL DEFAULT 0 COMMENT '更新日期',
) COMMENT '表'


