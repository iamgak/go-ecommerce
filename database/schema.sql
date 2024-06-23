
create database ecommerce;

CREATE TABLE seller_info (
  id SERIAL PRIMARY KEY ,
  user_id INTEGER UNIQUE NOT NULL,
  company_name VARCHAR(100) NOT NULL,
  pancard VARCHAR(50) UNIQUE NOT NULL,
  district_id INTEGER UNIQUE NOT NULL,
  state_id INTEGER UNIQUE NOT NULL,
  pincode INTEGER NOT NULL,
  addr TEXT NOT null
  );

CREATE TABLE loc_state (
  id SERIAL PRIMARY KEY,
  state_name VARCHAR(50) NOT NULL
)

CREATE TABLE loc_district (
  id SERIAL PRIMARY KEY,
  district_name VARCHAR(50) NOT NULL
)

-- Create users table 

CREATE TABLE user_listing (
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
  title VARCHAR(100) NOT NULL,
  quantity INTEGER NOT NULL,
  category_id INT NOT NULL,
  sub_category_id INT NOT NULL,
  descriptions text NOT NULL,
  product_addr_id INTEGER NOT NULL,
  price decimal(5,2) NOT NULL,
  user_id INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
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

CREATE TABLE cart (
  id SERIAL PRIMARY KEY ,
  product_id INT NOT NULL,
  user_id INTEGER NOT NULL,
  quantity INTEGER NOT NULL
);

CREATE TABLE order_listing (
  id SERIAL PRIMARY KEY ,
  product_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  quantity INTEGER NOT NULL,
  user_addr_id INTEGER DEFAULT NULL,
  payment_id INTEGER DEFAULT NULL,
  active BOOLEAN DEFAULT FALSE,
  shipping BOOLEAN DEFAULT FALSE,
  completed BOOLEAN DEFAULT FALSE,
  is_cancelled BOOLEAN DEFAULT FALSE,
  returning_id INTEGER DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

CREATE TABLE order_payment (
  id SERIAL PRIMARY KEY ,
  order_id INTEGER NOT NULL,
  amount DECIMAL(10,3) NOT NULL,
  done_by INTEGER NOT NULL,
  transaction_id VARCHAR(50) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

CREATE TABLE order_cancel (
  id SERIAL PRIMARY KEY ,
  order_id INTEGER NOT NULL,
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


CREATE TABLE order_shipping_addr (
  id SERIAL PRIMARY KEY ,
  order_id INTEGER NOT NULL,
  mobile VARCHAR(50) NOT NULL,
  region VARCHAR(50) NOT NULL,
  district VARCHAR(50) NOT NULL,
  pincode VARCHAR(50) NOT NULL,
  addr TEXT NOT NULL
);