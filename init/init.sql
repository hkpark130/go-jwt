SET NAMES 'utf8';

DROP TABLE IF EXISTS jwt_user;

CREATE TABLE jwt_user (
  id SERIAL PRIMARY KEY,
  email VARCHAR(200) NOT NULL,
  password VARCHAR(200) NOT NULL
);

INSERT INTO jwt_user (email, password) VALUES ('test@test.com', 'x');