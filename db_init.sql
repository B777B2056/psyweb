-- 新建数据库
DROP DATABASE IF EXISTS `psyweb`;
CREATE DATABASE `psyweb`;
USE `psyweb`;

-- 新建表
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
    `PhoneNumber` char(11),
    `VerificationCode` varchar(6),
    `SerialNumber` varchar(32),
    `Name` varchar(32),
    `Gender` tinyint(3) unsigned,
    `Age` tinyint(3) unsigned,
    `SASScore` float,
    `ESSScore` float,
    `ISIScore` float,
    `SDSScore` float,
    `Status` tinyint(3) unsigned,
    PRIMARY KEY (`PhoneNumber`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `staff`;
CREATE TABLE `staff` (
    `Id` varchar(16),
    `Password` varchar(16),
    PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 新建配置文件中记录的用户
USE mysql;
CREATE USER 'admin'@'%' identified by 'b777b2056';
GRANT ALL ON `psyweb`.* TO 'admin'@'%' identified by 'b777b2056' with grant option;
FLUSH PRIVILEGES;
