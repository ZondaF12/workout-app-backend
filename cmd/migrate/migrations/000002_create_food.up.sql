CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS foods (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  brand VARCHAR(255),
  calories INTEGER NOT NULL CHECK (calories >= 0),
  protein DECIMAL(8,2) NOT NULL CHECK (protein >= 0),
  carbs DECIMAL(8,2) NOT NULL CHECK (carbs >= 0),
  fat DECIMAL(8,2) NOT NULL CHECK (fat >= 0),
  fiber DECIMAL(8,2) NOT NULL CHECK (fiber >= 0),
  serving_size DECIMAL(8,2) NOT NULL CHECK (serving_size > 0),
  serving_unit VARCHAR(50) NOT NULL,
  verified BOOLEAN DEFAULT FALSE,
  user_id UUID REFERENCES users(id),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);