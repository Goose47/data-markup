CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   name VARCHAR(100) NOT NULL,
   email VARCHAR(255) UNIQUE NOT NULL,
   age INT CHECK (age > 0),
   created_at TIMESTAMP DEFAULT NOW()
);