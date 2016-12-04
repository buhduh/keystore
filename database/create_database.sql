/*
  TODO
  1. I will have to make a user and passwords, etc.
  2. Roll step 1 into the CI/CD pipeline
*/
drop table if exists users;

create table users (
  id mediumint not null auto_increment,
  name varchar(20) not null,
  password varchar(256) not null,
  salt varchar(32) not null,
  qr_secret varchar(32) not null, 
  primary key(id),
  constraint name unique(name)
) engine=innodb default charset=utf8;
