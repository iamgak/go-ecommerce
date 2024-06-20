
create database ecommerce;

CREATE TABLE seller_info (
  id SERIAL PRIMARY KEY ,
  user_id INTEGER UNIQUE NOT NULL,
  pancard VARCHAR(100) UNIQUE NOT NULL,
  addr text not null,
  state VARCHAR(50) not null,
  district VARCHAR(50) not null,
  pincode INTEGER NOT NULL
  );


-- Create users table 

CREATE TABLE users (
  id SERIAL PRIMARY KEY ,
  email VARCHAR(100) UNIQUE NOT NULL,
  password VARCHAR(100) NOT NULL,
  seller boolean DEFAULT FALSE,
  last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  login_token VARCHAR(100) DEFAULT NULL,
  active BOOLEAN DEFAULT FALSE,
  activation_token VARCHAR(100) DEFAULT NULL
);

--  Create forget_passw table 

CREATE TABLE user_forget_passw (
  id SERIAL PRIMARY KEY ,
  user_id INTEGER NOT NULL,
  uri VARCHAR(100) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  superseded BOOLEAN DEFAULT FALSE
);

--  Create user_log table 

CREATE TABLE user_log (
  id SERIAL PRIMARY KEY ,
  activity VARCHAR(50) NOT NULL,
  user_id INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  superseded BOOLEAN DEFAULT FALSE
);

-- Create Category table

CREATE TABLE category_main (
  id SERIAL PRIMARY KEY ,
  category VARCHAR(50) NOT NULL
);

-- Create Category table

CREATE TABLE category_sub (
  id SERIAL PRIMARY KEY ,
  sub_category VARCHAR(50) NOT NULL,
  category INTEGER NOT NULL
);

-- Create product table

CREATE TABLE product (
  id SERIAL PRIMARY KEY ,
  title VARCHAR(255) NOT NULL,
  quantity INTEGER NOT NULL,
  category INT NOT NULL,
  sub_category INT NOT NULL,
  descriptions text NOT NULL,
  price decimal(5,2) NOT NULL,
  user_id INTEGER NOT NULL,
  last_updated TIMESTAMP NOT NULL,
  is_deleted BOOLEAN DEFAULT FALSE
);

-- Create productFavourite table

CREATE TABLE favourite_product (
  id SERIAL PRIMARY KEY ,
  user_id INTEGER NOT NULL,
  product_id INTEGER NOT NULL
);

-- table product log

CREATE TABLE product_log (
  id SERIAL PRIMARY KEY ,
  activity VARCHAR(50) NOT NULL,
  user_id INTEGER NOT NULL,
  product_id INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  superseded BOOLEAN DEFAULT FALSE
);


-- order payment related work

CREATE TABLE order_payment (
  id SERIAL PRIMARY KEY ,
  total_amount decimal(10,2) not null,
  transaction_id INTEGER NOT NULL,
  transaction_status BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);


CREATE TABLE cart (
  id SERIAL PRIMARY KEY ,
  product_id INT NOT NULL,
  user_id INTEGER NOT NULL,
  quantity INTEGER NOT NULL,
  active BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE order_listing (
  id SERIAL PRIMARY KEY ,
  cart_id INTEGER UNIQUE NOT NULL,
  addr TEXT NOT NULL,
  total_price decimal(10,4) not null,
  payment_id INTEGER DEFAULT NULL,
  processed BOOLEAN DEFAULT FALSE,
  completed BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- order log

CREATE TABLE order_log (
  id SERIAL PRIMARY KEY ,
  activity VARCHAR(50) NOT NULL,
  order_id INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  superseded BOOLEAN DEFAULT FALSE
);