-- Active: 1714559811576@@127.0.0.1@3306@canteen-menu
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS menu_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    description TEXT,
    available BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    item_id INT NOT NULL,
    quantity INT NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    status ENUM('processing', 'ready') DEFAULT 'processing',
    order_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


    INSERT INTO menu_items (name, price, description, available) 
    VALUES 
        ('Chicken Burger', 5.99, 'Juicy chicken patty with lettuce and mayo', true),
        ('Vegetarian Pizza', 8.49, 'Delicious pizza topped with fresh vegetables', true),
        ('Pasta Carbonara', 7.99, 'Spaghetti with creamy carbonara sauce and bacon bits', true);

