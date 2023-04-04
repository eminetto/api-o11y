create database if not exists auth_db;
use auth_db;
create table if not exists user (id varchar(50),email varchar(255),password varchar(255),first_name varchar(100), last_name varchar(100), created_at datetime, updated_at datetime, PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
INSERT INTO user (id, email, password, first_name, last_name, created_at, updated_at)
SELECT 'adb8101e-cfe6-4a71-8594-ebc80af3a86d','eminetto@email.com',SHA1('12345'), 'Elton', 'Minetto', now(), null FROM DUAL
WHERE NOT EXISTS
    (SELECT email FROM user WHERE email='eminetto@email.com');