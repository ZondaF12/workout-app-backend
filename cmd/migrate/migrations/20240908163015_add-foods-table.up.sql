CREATE TABLE IF NOT EXISTS foods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    brand VARCHAR(100),
    default_serving_size DECIMAL(8, 2) NOT NULL,
    default_serving_unit VARCHAR(20) NOT NULL,
    calories INTEGER NOT NULL,
    protein DECIMAL(6, 2),
    carbs DECIMAL(6, 2),
    fat DECIMAL(6, 2),
    is_user_created BOOLEAN DEFAULT FALSE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes
CREATE INDEX idx_foods_name ON foods(name);
CREATE INDEX idx_foods_brand ON foods(brand);
CREATE INDEX idx_foods_is_user_created ON foods(is_user_created);
CREATE INDEX idx_foods_created_by ON foods(created_by);