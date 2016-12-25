drop table if exists passwords;
drop table if exists categories;
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

create table categories (
  id mediumint not null auto_increment,
  name varchar(20) not null,
  user_id mediumint not null,
  primary key(id),
  key user_id (user_id),
  foreign key (user_id) references users(id)
  on delete cascade,
  constraint name unique(name, user_id)
) engine=innodb default charset=utf8;

create table passwords (
  id mediumint not null auto_increment,
  password varchar(256) not null,
  user_name varchar(256),
  notes text,
  domain varchar(256),
  expires date,
  /*json column for allowable passwords, tbd*/
  rule_set text, 
  user_id mediumint not null,
  category_id mediumint not null,
  primary key(id),
  key user_id (user_id),
  key category_id (category_id),
  foreign key (user_id) references users(id)
  on delete cascade,
  foreign key (category_id) references categories(id)
  on delete cascade
) engine=innodb default charset=utf8;
