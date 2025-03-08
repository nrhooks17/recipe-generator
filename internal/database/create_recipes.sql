CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    recipe_name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT NULL,
    prep_time_minutes INT NULL,
    cook_time_minutes INT NULL,
    servings INT NULL,
    created_by INT REFERENCES users(id) NOT NULL,
    created_date DATE DEFAULT CURRENT_DATE NOT NULL,
    updated_by INT REFERENCES users(id) NOT NULL,
    updated_date DATE NOT NULL
)