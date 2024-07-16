# E-Commerce

Welcome to the E-Commerce! In this system there is three type of User Customer, Seller, Deliver & Admin.
This system allows users to buy product from Product listing, Add Orders in Cart, Make Payment using APIs,
or to sell product by Add/Remove Product, Add/Delete Discount, Change StockQuantity.

## Core Features

- User Authentication and Authorization:
 User registration and login
 Password recovery 
 Role-based access control (e.g. admin, seller, customer, deliver)
 
- Product Catalog:
 Product listing with images, descriptions, and prices
 Categories and subcategories
 Search functionality with filters and sorting options
 Product details page
 
- Shopping Cart:
 Add/remove products to/from the cart
 Update product quantities
 Display total cost
 Persist cart items across sessions
 
- Checkout Process:
 Address and shipping information
 <!-- Payment gateway integration (e.g., Stripe, PayPal) -->
 Order summary and confirmation
 <!-- Email notifications for order confirmations -->
 
- Order Management:
 Order history and tracking for customers
 <!-- Admin dashboard for managing orders -->
 Order status updates (e.g. pending, shipped, delivered)
 
- User Profile Management:
 Personal information management
 Address book for multiple shipping addresses
 Order history and tracking
<!--  
- Reviews and Ratings:
 Allow customers to leave reviews and ratings for products
 Display average ratings and reviews on product pages -->


## Technical and Infrastructure Requirements

- Backend Technologies:
Go (Golang) for building the backend services
RESTful APIs for client-server communication
Customizable Stateful token for secure authentication 
Mutex for locking goroutines in query and also in ratelimiter using package `mutex` and `sync`
Rate limiter to limit request using package `golang.org/x/time/rate` 
And added query timeout, db limit maxconnection, max idle/in-use connection using package `context`
Migrate go package for database versioning
Make for customize command

- Database:
PostGres for storing data
Redis for caching
<!-- 
- Payment Gateway Integration:
Stripe, PayPal, or other payment processors

- Email Service:
SMTP server or email service providers like SendGrid, TempMail or Mailgun for sending transactional emails

- Search Engine:
Elasticsearch for advanced search capabilities
 -->

## Security Measures

- Data Encryption:
HTTPS for secure communication
Encrypt sensitive data in the database

- Secure Coding Practices:
Input validation and sanitization
Protect against SQL injection, XSS, and CSRF attacks


## Scalability and Performance

- Database Optimization:
Optimize queries and indexing for performance

## Route Listing
- Home Page `https://localhost:4000/`
- User Account Activation `localhost:4000/api/user/activation_token/?activation_token=39c5f3bc4a65f44e625a88791a8440c63301b7f6s` 
- User Login `localhost:4000/api/user/login/`
- User Logout `localhost:4000/api/user/logout/`
- User Forget Password `localhost:4000/api/user/forget-password/`
- User Reset password `localhost:4000/api/user/new-password/?reset_token=4b39a1bf5da5b490217aeec50c392c57a08f6b33`

- User Page `https://localhost:4000/login/`, `https://localhost:4000/register/`, `https://localhost:4000/forget/password/`,`https://localhost:4000/user-info/`,`https://localhost:4000/update/user-info/`

- Add Product `https://localhost:4000/add/product/`,`https://localhost:4000/update/product/id/`, `https://localhost:4000/del/product/id/`, `https://localhost:4000/product/price/id/`, `https://localhost:4000/publish/id/`

- Listing `https://localhost:4000/listing/`,`https://localhost:4000/listing/with-search-pattern/`

- Order `https://localhost:4000/order/cart/id/`,`https://localhost:4000/order/user-info/id/`,
`https://localhost:4000/order/payment-info/id/`,`https://localhost:4000/order/submitted/id/`

## Conclusion
Contributions are welcome! If you'd like to contribute to this project, please fork the repository and submit a pull request with your changes.
Building an e-commerce platform involves a lot of components, and attention to detail is crucial for providing a seamless and secure shopping experience. This is just a minitare version of ecommerce without delivery management, smtp connection, without payment api and not account realtime goroutine, traffic, query optimisation.
