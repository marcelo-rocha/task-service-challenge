create table users(
  id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  login VARCHAR(40) NOT NULL,
  name VARCHAR(120) NOT NULL,
  kind ENUM('technician', 'manager') NOT NULL,
  active BOOLEAN NOT NULL DEFAULT true,
  UNIQUE KEY unique_login(login)
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



