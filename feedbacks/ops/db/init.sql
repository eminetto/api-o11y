create database if not exists feedbacks_db;
use feedbacks_db;
create table if not exists feedback (id varchar(50),email varchar(255),title varchar(255),body text,  created_at datetime, updated_at datetime, PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

