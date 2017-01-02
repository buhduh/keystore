/*
  @userName, @catFoo, @catBar, @catBaz, @addUserName
*/
insert users set name=@userName, password="foo", salt="foo", qr_secret="foo";
insert users set name=@addCatUserName, password="foo", salt="foo", qr_secret="foo";
insert categories(name, user_id) select @catFoo, id from users where name=@userName;
insert categories(name, user_id) select @catBar, id from users where name=@userName;
insert categories(name, user_id) select @catBaz, id from users where name=@userName;
