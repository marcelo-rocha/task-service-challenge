INSERT INTO users(login, name, kind) values('admin', 'administrator', 'manager');
INSERT INTO users(login, name, kind, manager_id) values('demo', 'demonstration', 'technician', LAST_INSERT_ID());

