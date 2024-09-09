CREATE TABLE meal_items (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    food_id INTEGER NOT NULL REFERENCES foods(id),
    date DATE NOT NULL,
    meal_type VARCHAR(20) NOT NULL,
    serving_size DECIMAL(8, 2) NOT NULL,
    serving_unit VARCHAR(20) NOT NULL,
    quantity DECIMAL(8, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CHECK (meal_type IN ('breakfast', 'lunch', 'dinner', 'snack'))
);

-- Create indexes
CREATE INDEX idx_meal_items_user_id_date ON meal_items(user_id, date);
CREATE INDEX idx_meal_items_food_id ON meal_items(food_id);
CREATE INDEX idx_meal_items_meal_type ON meal_items(meal_type);