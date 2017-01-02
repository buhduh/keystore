/*
  @userName, @catFoo, @catBar, @catBaz
  added categories from test should delete because of on delete cascade
*/
delete from users where name=@userName;
delete from users where name=@addCatUserName;
