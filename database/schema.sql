
create database ecommerce;

-- Create Category table
-- product category name

CREATE TABLE IF NOT EXISTS category_main (
  id BIGSERIAL PRIMARY KEY ,
  category VARCHAR(50) NOT NULL
);

-- Create Category table
-- product sub-category name

CREATE TABLE IF NOT EXISTS category_sub (
  id BIGSERIAL PRIMARY KEY ,
  sub_category VARCHAR(50) NOT NULL,
  category_id INTEGER NOT NULL
);

-- state name like KERLA,...

CREATE TABLE IF NOT EXISTS loc_state (
  id BIGSERIAL PRIMARY KEY,
  state_name VARCHAR(50) NOT NULL
);

-- payment types like CASHONDELIVERY, ONLINE, ...
CREATE TABLE IF NOT EXISTS type_payment (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL
);

-- district name like DISTRICT, ...
CREATE TABLE IF NOT EXISTS loc_district (
  id BIGSERIAL PRIMARY KEY,
  district_name VARCHAR(50) NOT NULL,
  state_id INTEGER NOT NULL
);

-- Create users table 

CREATE TABLE IF NOT EXISTS user_listing (
  id BIGSERIAL PRIMARY KEY ,
  email VARCHAR(100) UNIQUE NOT NULL,
  hashed_password BYTEA NOT NULL,
  last_login TIMESTAMP DEFAULT NULL,
  login_token VARCHAR(100) DEFAULT NULL,
  active BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  activation_token VARCHAR(100) DEFAULT NULL,
  version INTEGER NOT NULL DEFAULT 1

);


CREATE TABLE IF NOT EXISTS seller_listing (
  id BIGSERIAL PRIMARY KEY ,
  email VARCHAR(100) UNIQUE NOT NULL,
  hashed_password BYTEA NOT NULL,
  mobile VARCHAR(50) NOT NULL,
  company_name VARCHAR(50) NOT NULL,
  pancard VARCHAR(50) NOT NULL,
  region_id INTEGER NOT NULL,
  district_id INTEGER NOT NULL,
  pincode VARCHAR(50) NOT NULL,
  addr TEXT NOT NULL,  
  last_login TIMESTAMP DEFAULT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  login_token VARCHAR(100) DEFAULT NULL,
  active BOOLEAN DEFAULT FALSE,
  activation_token VARCHAR(100) DEFAULT NULL,
  version INTEGER NOT NULL DEFAULT 1

);

--  Create forget_passw table
CREATE TABLE IF NOT EXISTS user_forget_passw (
  id BIGSERIAL PRIMARY KEY ,
  user_id BIGINT NOT NULL,
  uri VARCHAR(100) NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  superseded BOOLEAN DEFAULT FALSE
);

--  Create user_log table
CREATE TABLE IF NOT EXISTS user_log (
  id BIGSERIAL PRIMARY KEY ,
  activity VARCHAR(50) NOT NULL,
  user_id BIGINT NOT NULL,
created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  superseded BOOLEAN DEFAULT FALSE
);

--  Create user_log table
CREATE TABLE IF NOT EXISTS seller_log (
  id BIGSERIAL PRIMARY KEY ,
  activity VARCHAR(50) NOT NULL,
  user_id BIGINT NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  superseded BOOLEAN DEFAULT FALSE
);

-- Create product table

CREATE TABLE IF NOT EXISTS product (
  id BIGSERIAL PRIMARY KEY ,
  title VARCHAR(100) NOT NULL,
  quantity INTEGER NOT NULL,
  category_id INT NOT NULL,
  sub_category_id INT NOT NULL,
  descriptions text NOT NULL,
  price decimal(10,2) NOT NULL,
  user_id BIGINT NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  active_from DATE DEFAULT NULL,
  active BOOLEAN DEFAULT TRUE,
  version INTEGER NOT NULL DEFAULT 1

);

-- Create productFavourite table

CREATE TABLE IF NOT EXISTS favourite_product (
  id BIGSERIAL PRIMARY KEY ,
  user_id BIGINT NOT NULL,
  product_id INTEGER NOT NULL
);

-- table product log

CREATE TABLE IF NOT EXISTS product_log (
  id BIGSERIAL PRIMARY KEY ,
  activity VARCHAR(50) NOT NULL,
  user_id BIGINT NOT NULL,
  product_id INTEGER NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  superseded BOOLEAN DEFAULT FALSE
);

-- order payment related work

CREATE TABLE IF NOT EXISTS cart (
  id BIGSERIAL PRIMARY KEY ,
  product_id INT NOT NULL,
  user_id BIGINT NOT NULL,
  quantity INTEGER NOT NULL,
  active BOOLEAN DEFAULT TRUE,
  version INTEGER NOT NULL DEFAULT 1

);

CREATE TABLE IF NOT EXISTS order_listing (
  id BIGSERIAL PRIMARY KEY ,
  cart_id INTEGER NOT NULL,
  product_id INTEGER NOT NULL,
  payment_method INTEGER NOT NULL,
  user_id BIGINT NOT NULL,
  addr_id INTEGER NOT NULL,
  price decimal(10,2) NOT NULL,
  quantity INTEGER NOT NULL,
  shipping_addr_id INTEGER DEFAULT NULL,
  payment_id INTEGER DEFAULT NULL,
  active BOOLEAN DEFAULT FALSE,
  dispatch BOOLEAN DEFAULT FALSE,
  is_cancelled BOOLEAN DEFAULT FALSE,
  returning_id INTEGER DEFAULT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  dispatch_at TIMESTAMP DEFAULT NULL,
  recieved_at TIMESTAMP DEFAULT NULL,
  version INTEGER NOT NULL DEFAULT 1

  );

CREATE TABLE IF NOT EXISTS order_payment (
  id BIGSERIAL PRIMARY KEY ,
  order_id INTEGER DEFAULT NULL,
  amount DECIMAL(10,3) NOT NULL,
  status BOOLEAN DEFAULT FALSE,
  transaction_id VARCHAR(50) NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  version INTEGER NOT NULL DEFAULT 1
  );

CREATE TABLE IF NOT EXISTS order_cancel (
  id BIGSERIAL PRIMARY KEY ,
  order_id INTEGER NOT NULL,
  delivery_recieved BOOLEAN DEFAULT FALSE,
  seller_recieved BOOLEAN DEFAULT FALSE,
  bad_condition BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  version INTEGER NOT NULL DEFAULT 1
);


-- order log

CREATE TABLE IF NOT EXISTS order_log (
  id BIGSERIAL PRIMARY KEY ,
  activity VARCHAR(50) NOT NULL,
  order_id INTEGER NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  superseded BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS user_addr (
  id BIGSERIAL PRIMARY KEY ,
  user_id BIGINT NOT NULL,
  mobile VARCHAR(50) NOT NULL,
  region_id VARCHAR(50) NOT NULL,
  district_id VARCHAR(50) NOT NULL,
  pincode VARCHAR(50) NOT NULL,
  addr TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS product_origin_addr (
  id BIGSERIAL PRIMARY KEY ,
  order_id INTEGER NOT NULL,
  mobile VARCHAR(50) NOT NULL,
  region_id INTEGER NOT NULL,
  district_id INTEGER NOT NULL,
  pincode VARCHAR(50) NOT NULL,
  addr TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  name text NOT NULL,
  email citext UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  activated bool NOT NULL,
  version integer NOT NULL DEFAULT 1
);

// for login_token session 
CREATE TABLE IF NOT EXISTS tokens (
hash bytea PRIMARY KEY,
user_id bigint NOT NULL REFERENCES user_listing ON DELETE CASCADE,
expiry timestamp(0) with time zone NOT NULL,
scope text NOT NULL
);
