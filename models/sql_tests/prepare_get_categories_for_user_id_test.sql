insert users set name="foo", password="foo", salt="foo", qr_secret="foo";
insert categories(name, user_id) select "category_foo", id from users where name="foo";
insert categories(name, user_id) select "category_bar", id from users where name="foo";
insert categories(name, user_id) select "category_baz", id from users where name="foo";
