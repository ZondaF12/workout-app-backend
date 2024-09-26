CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  first_name varchar(255) NOT NULL,
  last_name varchar(255) NOT NULL,
  email citext UNIQUE NOT NULL,
  username varchar(255) UNIQUE NOT NULL,
  password bytea NOT NULL,
  bio text,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);