-- Create the transactions table
CREATE TABLE transactions (
                              transaction_id SERIAL PRIMARY KEY,
                              user_id INT NOT NULL REFERENCES users(user_id),
                              cart_id INT NOT NULL REFERENCES carts(cart_id),
                              status VARCHAR(50) DEFAULT 'Pending Payment',
                              created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);