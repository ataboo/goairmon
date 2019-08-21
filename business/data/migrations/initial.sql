CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT,
  passwordhash TEXT,
  lastlogin DATETIME
);
