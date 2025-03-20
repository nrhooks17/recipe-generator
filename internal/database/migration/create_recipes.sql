
DO $$
-- start muh transackshaaan
BEGIN

INSERT INTO recipes (recipe_name, description, prep_time_minutes,  cook_time_minutes, servings, created_by, created_date, updated_by, updated_date)
VALUES ('Pancakes', 'Delicious pancakes', 15, 10, 4, 1, NOW(), 1, NOW()),
('Salad', 'Healthy salad', 10, 5, 2, 1, NOW(), 1, NOW()),
('Pizza', 'Tasty pizza', 20, 15, 6, 1, NOW(), 1, NOW()),
('Sushi', 'Delicious sushi', 30, 20, 4, 1, NOW(), 1, NOW()),
('Tacos', 'Delicious tacos', 15, 10, 4, 1, NOW(), 1, NOW());

-- ingredients insertion
INSERT INTO ingredients(recipe_id, unit_of_measurement, unit_amount, ingredient_name, created_by, created_date, updated_by, updated_date)
Values((SELECT id FROM recipes WHERE recipe_name = 'Pancakes'), 'cup', 2, 'flour', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pancakes'), 'cup', 1, 'sugar', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pancakes'), 'cup', 1, 'eggs', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pancakes'), 'cup', 1, 'milk', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pancakes'), 'cup', 1, 'butter', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pancakes'), 'cup', 1, 'vanilla extract', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pancakes'), 'cup', 1, 'salt', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pancakes'), 'cup', 1, 'baking powder', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pancakes'), 'cup', 1, 'baking soda', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Salad'), 'cup', 1, 'lettuce', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Salad'), 'cup', 1, 'tomato', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Salad'), 'cup', 1, 'cucumber', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Salad'), 'cup', 1, 'onion', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Salad'), 'cup', 1, 'olive oil', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Salad'), 'cup', 1, 'salt', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Salad'), 'cup', 1, 'pepper', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Salad'), 'cup', 1, 'mayonnaise', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Pizza'), 'cup', 1, 'flour', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Pizza'), 'cup', 1, 'yeast', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Pizza'), 'cup', 1, 'salt', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Pizza'), 'cup', 1, 'tomato sauce', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Pizza'), 'cup', 1, 'mozzarella cheese', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Pizza'), 'cup', 1, 'pepperoni', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Pizza'), 'cup', 1, 'olive oil', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Sushi'), 'cup', 1, 'rice', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Sushi'), 'cup', 1, 'seaweed', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Sushi'), 'cup', 1, 'tuna', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Sushi'), 'cup', 1, 'avocado', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Sushi'), 'cup', 1, 'sushi rice', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Sushi'), 'cup', 1, 'soy sauce', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Sushi'), 'cup', 1, 'sesame oil', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Tacos'), 'cup', 1, 'ground beef', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Tacos'), 'cup', 1, 'lettuce', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Tacos'), 'cup', 1, 'tomato', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Tacos'), 'cup', 1, 'salsa', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Tacos'), 'cup', 1, 'sour cream', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Tacos'), 'cup', 1, 'cilantro', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes where recipe_name = 'Tacos'), 'cup', 1, 'taco shells', 1, NOW(), 1, NOW());

INSERT INTO procedure_steps (recipe_id, step, created_by, created_date, updated_by, updated_date)
VALUES ((SELECT id FROM recipes WHERE recipe_name = 'Pancakes'), 'First, make sure that the oven is on.', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pancakes'), 'Next, mix all the ingredients together.', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pancakes'), 'Then, pour the batter into a pan and cook until it is done.', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pancakes'), 'Finally, enjoy your delicious pancakes!', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Salad'), 'First, cut the lettuce, tomatoes and onions into small pieces.', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Salad'), 'Next, mix all the ingredients together.', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Salad'), 'Then, pour the salad into a bowl and enjoy your delicious salad!', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pizza'), 'First, mix all the ingredients together.', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pizza'), 'Next, pour the pizza into a pan and cook until it is done.', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Pizza'), 'Then, enjoy your delicious pizza!', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Sushi'), 'First, mix all the ingredients together.', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Sushi'), 'Next, pour the sushi into a pan and cook until it is done.', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Sushi'), 'Then, enjoy your delicious sushi!', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Tacos'), 'First, mix all the ingredients together.', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Tacos'), 'Next, pour the tacos into a pan and cook until it is done.', 1, NOW(), 1, NOW()),
((SELECT id FROM recipes WHERE recipe_name = 'Tacos'), 'Then, enjoy your delicious tacos!', 1, NOW(), 1, NOW());

-- end muh transackshaaan

EXCEPTION WHEN OTHERS THEN
    -- roll back when any errors occur.

    RAISE NOTICE 'an error occurred: %', SQLERRM;
    ROLLBACK;
END$$;

