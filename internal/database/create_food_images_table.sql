CREATE TABLE food_images (
    id SERIAL PRIMARY KEY NOT NULL,
    recipe_id INT REFERENCES recipes(id) NOT NULL,
    filepath VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    size_bytes BIGINT NOT NULL,
    created_by INT REFERENCES users(id) NOT NULL,
    created_date DATE DEFAULT CURRENT_DATE NOT NULL,
    updated_by INT REFERENCES users(id) NOT NULL,
    updated_date DATE NOT NULL
)