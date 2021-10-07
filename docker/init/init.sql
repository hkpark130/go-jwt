SET NAMES 'utf8';

DROP TABLE IF EXISTS jwt_users;

CREATE TABLE jwt_users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(200) NOT NULL,
  password VARCHAR(200) NOT NULL,
  created_at DATE,
  updated_at DATE,
  deleted_at DATE
);

INSERT INTO jwt_users (email, password) VALUES ('test@test.com', 'x');
