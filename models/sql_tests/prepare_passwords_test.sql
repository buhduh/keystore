/*
  required vars: userFoo, catFoo, catBar, passFoo, passBar, 
  userUpdate, passUpdate, notesUpdate, domainUpdate, expiresUpdate,
  rulesUpdate, oldCatUpdate, newCatUpdate, nameUpdate
*/

insert users set name=@userFoo, password="foo", salt="foo", qr_secret="foo";
insert categories(name, user_id) select @catFoo, id from users where name=@userFoo;
insert categories(name, user_id) select @catBar, id from users where name=@userFoo;

insert passwords(password, user_name, notes, domain, expires, rule_set, user_id, category_id)
  select @passFoo, "user_name", "notes", "domain", "2001-01-24", "some rules", u.id, c.id 
    from
      users u join categories c on u.id=c.user_id 
    where 
      u.name=@userFoo and c.name=@catFoo;

insert passwords(password, user_name, notes, domain, expires, rule_set, user_id, category_id)
  select @passBar, "user_name", "notes", "domain", "2001-01-24", "some rules", u.id, c.id 
    from
      users u join categories c on u.id=c.user_id 
    where 
      u.name=@userFoo and c.name=@catBar;

/*Update test*/
insert users set name=@userUpdate, password="foo", salt="foo", qr_secret="foo";
insert categories(name, user_id) select @oldCatUpdate, id from users where name=@userUpdate;
insert categories(name, user_id) select @newCatUpdate, id from users where name=@userUpdate;

insert passwords(password, user_name, notes, domain, expires, rule_set, user_id, category_id)
  select @passUpdate, @nameUpdate, @notesUpdate, @domainUpdate, @expiresUpdate, @rulesUpdate, u.id, c.id 
    from
      users u join categories c on u.id=c.user_id 
    where 
      u.name=@userUpdate and c.name=@oldCatUpdate;
