create database if not exists votes_db;
use votes_db;
create table if not exists vote (id varchar(50),email varchar(255),talk_name varchar(255), score int,  created_at datetime, updated_at datetime, PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
