-- init-schema.sql

-- Create the 'users' table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    mobile VARCHAR(50) NOT NULL,
    user_name VARCHAR(50) NOT NULL
);

-- You can add other tables or initial data here