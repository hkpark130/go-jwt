SET NAMES 'utf8';

DROP TABLE IF EXISTS jwt_users;

CREATE TABLE jwt_users (
  id SERIAL NOT NULL PRIMARY KEY,
  email VARCHAR(200) NOT NULL DEFAULT 'user@example.com',
  password VARCHAR(200) NOT NULL DEFAULT 'password',
  created_at DATE DEFAULT CURRENT_TIMESTAMP,
  updated_at DATE,
  deleted_at DATE
);

INSERT INTO jwt_users (email, password) VALUES ('test@test.com', 'x');
