create database if not exists auth_db;
use auth_db;
create table if not exists user (id varchar(50),email varchar(255),password varchar(255),first_name varchar(100), last_name varchar(100), created_at datetime, updated_at datetime, PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
