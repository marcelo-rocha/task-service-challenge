INSERT INTO users(login, name, kind) values('admin', 'administrator', 'manager');
SET @adminId = LAST_INSERT_ID();
INSERT INTO users(login, name, kind, manager_id) values('demo', 'demonstration', 'technician', @adminId);
INSERT INTO users(login, name, kind, manager_id) values('operator', 'operator demonstration', 'technician', @adminId);

