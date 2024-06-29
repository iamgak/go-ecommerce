
create database ecommerce;

-- Create Category table
-- product category name

CREATE TABLE category_main (
  id SERIAL PRIMARY KEY ,
  category VARCHAR(50) NOT NULL
);

-- Create Category table
-- product sub-category name

CREATE TABLE category_sub (
  id SERIAL PRIMARY KEY ,
  sub_category VARCHAR(50) NOT NULL,
  category_id INTEGER NOT NULL
);

-- state name like KERLA,...

CREATE TABLE loc_state (
  id SERIAL PRIMARY KEY,
  state_name VARCHAR(50) NOT NULL
);

-- payment types like CASHONDELIVERY, ONLINE, ...
CREATE TABLE type_payment (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL
);

-- district name like DISTRICT, ...
CREATE TABLE loc_district (
  id SERIAL PRIMARY KEY,
  district_name VARCHAR(50) NOT NULL,
  state_id INTEGER NOT NULL
);

-- Create users table 

CREATE TABLE user_listing (
  id SERIAL PRIMARY KEY ,
  email VARCHAR(100) UNIQUE NOT NULL,
  hashed_password VARCHAR(100) NOT NULL,
  last_login TIMESTAMP DEFAULT NULL,
  login_token VARCHAR(100) DEFAULT NULL,
  active BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  activation_token VARCHAR(100) DEFAULT NULL
);


CREATE TABLE seller_listing (
  id SERIAL PRIMARY KEY ,
  email VARCHAR(100) UNIQUE NOT NULL,
  hashed_password VARCHAR(100) NOT NULL,
  mobile VARCHAR(50) NOT NULL,
  company_name VARCHAR(50) NOT NULL,
  pancard VARCHAR(50) NOT NULL,
  region_id INTEGER NOT NULL,
  district_id INTEGER NOT NULL,
  pincode VARCHAR(50) NOT NULL,
  addr TEXT NOT NULL,  
  last_login TIMESTAMP DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
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

--  Create user_log table
CREATE TABLE seller_log (
  id SERIAL PRIMARY KEY ,
  activity VARCHAR(50) NOT NULL,
  user_id INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  superseded BOOLEAN DEFAULT FALSE
);

-- Create product table

CREATE TABLE product (
  id SERIAL PRIMARY KEY ,
  title VARCHAR(100) NOT NULL,
  quantity INTEGER NOT NULL,
  category_id INT NOT NULL,
  sub_category_id INT NOT NULL,
  descriptions text NOT NULL,
  price decimal(10,2) NOT NULL,
  user_id INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  active_from DATE DEFAULT NULL,
  active BOOLEAN DEFAULT TRUE
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

CREATE TABLE cart (
  id SERIAL PRIMARY KEY ,
  product_id INT NOT NULL,
  user_id INTEGER NOT NULL,
  quantity INTEGER NOT NULL,
  active BOOLEAN DEFAULT TRUE
);

CREATE TABLE order_listing (
  id SERIAL PRIMARY KEY ,
  cart_id INTEGER NOT NULL,
  product_id INTEGER NOT NULL,
  payment_method INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  addr_id INTEGER NOT NULL,
  price decimal(10,2) NOT NULL,
  quantity INTEGER NOT NULL,
  shipping_addr_id INTEGER DEFAULT NULL,
  payment_id INTEGER DEFAULT NULL,
  active BOOLEAN DEFAULT FALSE,
  dispatch BOOLEAN DEFAULT FALSE,
  is_cancelled BOOLEAN DEFAULT FALSE,
  returning_id INTEGER DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  dispatch_at TIMESTAMP DEFAULT NULL,
  recieved_at TIMESTAMP DEFAULT NULL
  );

CREATE TABLE order_payment (
  id SERIAL PRIMARY KEY ,
  order_id INTEGER DEFAULT NULL,
  amount DECIMAL(10,3) NOT NULL,
  status BOOLEAN DEFAULT FALSE,
  transaction_id VARCHAR(50) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

CREATE TABLE order_cancel (
  id SERIAL PRIMARY KEY ,
  order_id INTEGER NOT NULL,
  delivery_recieved BOOLEAN DEFAULT FALSE,
  seller_recieved BOOLEAN DEFAULT FALSE,
  bad_condition BOOLEAN DEFAULT FALSE,
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

CREATE TABLE user_addr (
  id SERIAL PRIMARY KEY ,
  user_id INTEGER NOT NULL,
  mobile VARCHAR(50) NOT NULL,
  region_id VARCHAR(50) NOT NULL,
  district_id VARCHAR(50) NOT NULL,
  pincode VARCHAR(50) NOT NULL,
  addr TEXT NOT NULL
);

CREATE TABLE product_origin_addr (
  id SERIAL PRIMARY KEY ,
  order_id INTEGER NOT NULL,
  mobile VARCHAR(50) NOT NULL,
  region_id INTEGER NOT NULL,
  district_id INTEGER NOT NULL,
  pincode VARCHAR(50) NOT NULL,
  addr TEXT NOT NULL
);