CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE,
  password TEXT
);

CREATE TABLE sessions (
  id TEXT PRIMARY KEY,
  user_id INT REFERENCES users(id)
);

CREATE TABLE comments (
  id SERIAL PRIMARY KEY,
  author TEXT,
  content TEXT
);

-- Insert a dummy user and a normal comment
INSERT INTO users (username, password) VALUES ('admin', 'password123');
INSERT INTO comments (author, content) VALUES ('System', 'Welcome to the secure comment board!');