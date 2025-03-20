CREATE TABLE ingredients (
    id SERIAL PRIMARY KEY,
    recipe_id INT REFERENCES recipes(id) NOT NULL,
    unit_of_measurement VARCHAR(255) NOT NULL,
    unit_amount DOUBLE PRECISION NOT NULL,
    ingredient_name VARCHAR(255) NOT NULL,
    created_by INT REFERENCES users(id) NOT NULL,
    created_date DATE DEFAULT CURRENT_DATE NOT NULL,
    updated_by INT REFERENCES users(id) NOT NULL,
    updated_date DATE NOT NULL
)