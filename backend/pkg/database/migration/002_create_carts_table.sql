-- Create the carts table
CREATE TABLE carts (
                       cart_id SERIAL PRIMARY KEY,
                       user_id INT NOT NULL REFERENCES users(user_id),
                       product_id INT NOT NULL,
                       quantity INT NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);