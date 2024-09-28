CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS roles (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO
    roles (name, level, description)
VALUES
    ('admin', 3, 'The highest level of access. Can do anything.'),
    ('subscriber', 2, 'A user who has subscribed to the service.'),
    ('user', 1, 'The default role for all users.');