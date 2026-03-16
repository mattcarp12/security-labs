-- init.sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE,
  password TEXT
);

CREATE TABLE sessions (
  id TEXT PRIMARY KEY,
  user_id INT REFERENCES users(id),
  csrf_token TEXT
);

CREATE TABLE transfers (
  id SERIAL PRIMARY KEY,
  user_id INT,
  amount INT
);

-- Insert a dummy user so we can test the login later
INSERT INTO users (username, password) VALUES ('admin', 'password123');