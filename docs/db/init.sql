-- DDL
DROP DATABASE IF EXISTS `tester_example_tasks_db`;
CREATE DATABASE `tester_example_tasks_db`;
USE `tester_example_tasks_db`;

CREATE TABLE `tasks` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `title` varchar(255) NOT NULL,
    `description` text NOT NULL,
    `done` boolean NOT NULL DEFAULT false,
    PRIMARY KEY (`id`)
);
