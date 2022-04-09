SET NAMES 'utf8';

DROP TABLE IF EXISTS jwt_users;

CREATE TABLE jwt_users (
  id SERIAL NOT NULL PRIMARY KEY,
  email VARCHAR(200) NOT NULL DEFAULT 'user@example.com',
  permission VARCHAR(200) NOT NULL DEFAULT 'role_user',
  password VARCHAR(200) NOT NULL DEFAULT 'password',
  created_at DATE DEFAULT CURRENT_TIMESTAMP,
  updated_at DATE,
  deleted_at DATE
);

INSERT INTO jwt_users (email, permission, password) VALUES ('admin@test.com', 'role_admin', '$2a$10$fkUtvFeOb17E7fF0tNV1tOtyZHqPp1IDDOvvXs9SxuTceGlU5lmiu');
INSERT INTO jwt_users (email, permission, password) VALUES ('test@test.com', 'role_user', '$2a$10$fkUtvFeOb17E7fF0tNV1tOtyZHqPp1IDDOvvXs9SxuTceGlU5lmiu');
