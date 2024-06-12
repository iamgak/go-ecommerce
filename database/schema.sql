-- Create database if it doesn't exist
CREATE DATABASE IF NOT EXISTS `ecommerce`;
USE `ecommerce`;

-- Create Category table
CREATE TABLE IF NOT EXISTS `category_main` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `category` VARCHAR(50) NOT NULL
);

-- Create Category table
CREATE TABLE IF NOT EXISTS `category_sub` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `sub-category` VARCHAR(50) NOT NULL,
  `category` INT(11) NOT NULL
);

-- Create UserFavouriteListing table
CREATE TABLE IF NOT EXISTS `object_fav` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `uid` int(11) NOT NULL,
  `object_id` int(11) NOT NULL
);

-- Create users table 

CREATE TABLE IF NOT EXISTS `seller` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `email` VARCHAR(100) UNIQUE NOT NULL,
  `password` VARCHAR(100) NOT NULL,
  `last_login` DATETIME DEFAULT current_timestamp(),
  `login_token` VARCHAR(100) DEFAULT NULL,
  `active` tinyint(1) DEFAULT 0,
  `pancard` VARCHAR(10) NOT NULL,
  `activation_token` VARCHAR(100) DEFAULT NULL
);


-- Create users table 

CREATE TABLE IF NOT EXISTS `user` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `email` VARCHAR(100) UNIQUE NOT NULL,
  `password` VARCHAR(100) NOT NULL,
  `last_login` DATETIME DEFAULT current_timestamp(),
  `login_token` VARCHAR(100) DEFAULT NULL,
  `active` tinyint(1) DEFAULT 0,
  `activation_token` VARCHAR(100) DEFAULT NULL
);

--  Create forget_passw table 

CREATE TABLE IF NOT EXISTS `user_forget_passw` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `uid` int(11) NOT NULL,
  `uri` VARCHAR(100) NOT NULL,
  `created_at` DATETIME DEFAULT current_timestamp(),
  `superseded` tinyint(1) DEFAULT 0
);

--  Create user_log table 

CREATE TABLE IF NOT EXISTS `user_log` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `activity` VARCHAR(50) NOT NULL,
  `uid` int(11) NOT NULL,
  `created_at` DATETIME DEFAULT current_timestamp(),
  `superseded` tinyint(1) DEFAULT 0
);

-- Create object table
CREATE TABLE IF NOT EXISTS `object` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `quantity` int(11) NOT NULL,
  `category` INT NOT NULL,
  `sub-category` INT NOT NULL,
  `descriptions` text NOT NULL,
  `price` decimal(5,2) NOT NULL,
  `uid` int(11) NOT NULL,
  `last_updated` DATETIME NOT NULL,
  `is_deleted` tinyint(1) DEFAULT 0
);

-- Create order-status table
CREATE TABLE IF NOT EXISTS `order_status` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `tab1` tinyint(1) DEFAULT 0,
  `tab2` tinyint(1) DEFAULT 0,
  `tab3` tinyint(1) DEFAULT 0,
  `tab4` tinyint(1) DEFAULT 0,
  `tab5` tinyint(1) DEFAULT 0,
  `final` tinyint(1) DEFAULT 0,
  `order_id` int(11) NOT NULL
);


CREATE TABLE IF NOT EXISTS `order_listing` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `object_id` int(11) NOT NULL,
  `uid` int(11) NOT NULL,
  `quantity` int(11) NOT NULL,
  `price` decimal(5,2) NOT NULL,
  `is_active` tinyint(1) DEFAULT 0,
  `completed` tinyint(1) DEFAULT 0
);

CREATE TABLE IF NOT EXISTS `order_contact` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `order_id` int(11) NOT NULL,
  `addr` VARCHAR(50) NOT NULL,
  `phone` VARCHAR(50) NOT NULL
);


CREATE TABLE IF NOT EXISTS `order_payment` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `order_id` int(11) NOT NULL,
  `payment_id` int(11) NOT NULL
);

CREATE TABLE IF NOT EXISTS `object_log` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `activity` VARCHAR(50) NOT NULL,
  `uid` int(11) NOT NULL,
  `object_id` int(11) NOT NULL,
  `created_at` DATETIME DEFAULT current_timestamp(),
  `superseded` tinyint(1) DEFAULT 0
);


CREATE TABLE IF NOT EXISTS `order_log` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `activity` VARCHAR(50) NOT NULL,
  `order_id` int(11) NOT NULL,
  `created_at` DATETIME DEFAULT current_timestamp(),
  `superseded` tinyint(1) DEFAULT 0
);
