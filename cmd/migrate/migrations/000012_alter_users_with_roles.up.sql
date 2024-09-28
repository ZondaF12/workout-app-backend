CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Ensure the roles table has the 'user' role
INSERT INTO roles (id, name)
VALUES ('0262b02f-21b6-4e07-a805-d3b634f9780c', 'user')
ON CONFLICT (name) DO NOTHING;

ALTER TABLE users
ADD COLUMN role_id uuid REFERENCES roles(id);

UPDATE users
SET role_id = (
  SELECT id
  FROM roles
  WHERE name = 'user'
);

ALTER TABLE users
ALTER COLUMN role_id SET NOT NULL;