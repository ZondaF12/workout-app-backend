CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Index on name for fast text searches
CREATE INDEX IF NOT EXISTS idx_foods_name ON foods USING GIN (to_tsvector('english', name));
CREATE INDEX IF NOT EXISTS idx_foods_name_brand ON foods (name, brand);
CREATE INDEX IF NOT EXISTS idx_foods_user_added ON foods (user_id);

CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

CREATE INDEX IF NOT EXISTS idx_meals_user_id_date ON meals (user_id, date);

CREATE INDEX IF NOT EXISTS idx_meal_entries_meal_id ON meal_entries (meal_id);
CREATE INDEX IF NOT EXISTS idx_meal_entries_food_id ON meal_entries (food_id);
CREATE INDEX IF NOT EXISTS idx_meal_entries_consumed_at ON meal_entries (consumed_at);