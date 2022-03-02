create table users(
  id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  login VARCHAR(40) NOT NULL,
  name VARCHAR(120) NOT NULL,
  kind ENUM('technician', 'manager') NOT NULL,
  active BOOLEAN NOT NULL DEFAULT true,
  manager_id int, 
  UNIQUE KEY unique_login(login),
  CONSTRAINT sr_fk_emp_man FOREIGN KEY (manager_id) REFERENCES users(id) 
);

create table tasks(
  id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT, 
  name VARCHAR(100) NOT NULL,
  summary VARCHAR(2500) NULL,
  creation_date TIMESTAMP NOT NULL DEFAULT now(),
  finish_date TIMESTAMP NULL,
  user_id BIGINT NOT NULL,
  INDEX user_idx(user_id),
  FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);



