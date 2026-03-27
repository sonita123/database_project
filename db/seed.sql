-- =========================================================
-- USERS
-- =========================================================
INSERT INTO users (first_name,last_name,email,balance,user_type) VALUES
                                                                     ('Ali','Ahmadi','ali.ahmadi@email.com',500000,'regular'),
                                                                     ('Sara','Hosseini','sara.hosseini@email.com',250000,'vip'),
                                                                     ('Reza','Karimi','reza.karimi@email.com',180000,'regular'),
                                                                     ('Maryam','Mohammadi','maryam.mohammadi@email.com',320000,'regular'),
                                                                     ('Hossein','Rezaei','hossein.rezaei@email.com',750000,'vip'),
                                                                     ('Fateme','Moradi','fateme.moradi@email.com',120000,'regular'),
                                                                     ('Mohammad','Safari','mohammad.safari@email.com',430000,'regular'),
                                                                     ('Zahra','Jafari','zahra.jafari@email.com',290000,'regular'),
                                                                     ('Amir','Abbasi','amir.abbasi@email.com',610000,'regular'),
                                                                     ('Narges','Tavakoli','narges.tavakoli@email.com',95000,'regular');

-- =========================================================
-- SUPPORTERS
-- =========================================================
INSERT INTO supporters (first_name,last_name,image_url,email) VALUES
                                                                  ('Kamran','Shirazi','https://example.com/kamran.jpg','kamran@example.com'),
                                                                  ('Leila','Bagheri','https://example.com/leila.jpg','leila@example.com'),
                                                                  ('Davood','Sadeghi','https://example.com/davood.jpg','davood@example.com');

-- =========================================================
-- ADMINS
-- =========================================================
INSERT INTO admins (username,password) VALUES
                                           ('admin1','AdminPass123!'),
                                           ('admin2','SecurePass456!');

-- =========================================================
-- SESSIONS
-- =========================================================
INSERT INTO sessions (token,user_id,role,expires_at) VALUES
                                                         ('TOKEN_1_64CHAR_LONG_STRING_________0000000000000000000000',1,'user',DATEADD(HOUR,3,GETDATE())),
                                                         ('TOKEN_2_64CHAR_LONG_STRING_________0000000000000000000000',2,'user',DATEADD(HOUR,3,GETDATE())),
                                                         ('TOKEN_3_64CHAR_LONG_STRING_________0000000000000000000000',3,'supporter',DATEADD(HOUR,3,GETDATE()));

-- =========================================================
-- VIP USERS
-- =========================================================
INSERT INTO vip_users (user_id,start_date,end_date) VALUES
                                                        (2,DATEADD(DAY,-30,GETDATE()),DATEADD(DAY,335,GETDATE())),
                                                        (5,DATEADD(DAY,-10,GETDATE()),DATEADD(DAY,355,GETDATE()));

-- =========================================================
-- ADDRESSES
-- =========================================================
INSERT INTO addresses (user_id,city,street,postal_code,is_default) VALUES
                                                                       (1,'Tehran','Valiasr St 12','123456',1),
                                                                       (2,'Isfahan','Chahar Bagh 7','654321',1),
                                                                       (3,'Mashhad','Imam Reza 33','987654',1),
                                                                       (4,'Tehran','Azadi 88','111222',1),
                                                                       (5,'Shiraz','Zand Blvd 14','333444',1);

-- =========================================================
-- SELLERS
-- =========================================================
INSERT INTO sellers (user_id) VALUES
                                  (1),(3),(5),(7),(9);

-- =========================================================
-- STALLS
-- =========================================================
INSERT INTO stalls (seller_id,name,status,approved_by) VALUES
                                                           (1,'Ali Tech Shop','active',1),
                                                           (2,'Reza Book Store','active',2),
                                                           (3,'VIP Electronics','active',1),
                                                           (4,'Mohammad Fashion','active',3),
                                                           (5,'Amir Sports','active',2);

-- =========================================================
-- PRODUCTS
-- =========================================================
INSERT INTO products (stall_id,name,price,stock) VALUES
                                                     (1,'Laptop Pro 15',35000000,10),
                                                     (1,'Wireless Mouse',450000,50),
                                                     (1,'USB-C Hub',780000,30),
                                                     (2,'Clean Code',320000,25),
                                                     (2,'Design Patterns',280000,20),
                                                     (2,'Pragmatic Programmer',350000,15),
                                                     (3,'iPhone 15',55000000,8),
                                                     (3,'Samsung Galaxy S24',48000000,12),
                                                     (3,'AirPods Pro',9500000,20),
                                                     (4,'Linen Shirt',1200000,40),
                                                     (4,'Denim Jacket',2800000,25),
                                                     (4,'Running Shoes',3500000,18),
                                                     (5,'Yoga Mat',850000,35),
                                                     (5,'Dumbbell Set',4200000,10);

-- =========================================================
-- PRODUCT VIEWS
-- =========================================================
INSERT INTO product_views (user_id,product_id,viewed_at) VALUES
                                                             (1,1,DATEADD(DAY,-2,GETDATE())),
                                                             (1,2,DATEADD(DAY,-1,GETDATE())),
                                                             (2,7,DATEADD(DAY,-5,GETDATE())),
                                                             (3,4,DATEADD(DAY,-6,GETDATE())),
                                                             (4,10,DATEADD(DAY,-1,GETDATE()));

-- =========================================================
-- CARTS
-- =========================================================
INSERT INTO carts (user_id,locked) VALUES
                                       (1,0),(2,0),(3,0),(4,0),(5,0);

-- =========================================================
-- CART ITEMS
-- =========================================================
INSERT INTO cart_items (cart_id,product_id,quantity) VALUES
                                                         (1,2,1),(1,3,2),(2,7,1),(3,4,1),(4,10,1),(5,6,1);

-- =========================================================
-- ORDERS
-- =========================================================
INSERT INTO orders (user_id,cart_id,address_id,total_price,status) VALUES
                                                                       (1,1,1,600000,'pending'),
                                                                       (2,2,2,55000000,'confirmed'),
                                                                       (3,3,3,320000,'shipped'),
                                                                       (4,4,4,9500000,'delivered'),
                                                                       (5,5,5,350000,'cancelled');

-- =========================================================
-- ORDER ITEMS
-- =========================================================
INSERT INTO order_items (order_id,product_id,quantity,price) VALUES
                                                                 (1,4,2,320000),
                                                                 (1,5,1,280000),
                                                                 (2,7,1,55000000),
                                                                 (3,4,1,320000),
                                                                 (4,9,1,9500000),
                                                                 (5,6,1,350000);

-- =========================================================
-- DISCOUNT CODES
-- =========================================================
INSERT INTO discount_codes (code,discount_type,percentage,expiration_date,supporter_id) VALUES
                                                                                            ('SAVE10','percentage',10,DATEADD(DAY,30,GETDATE()),1),
                                                                                            ('SAVE20','percentage',20,DATEADD(DAY,15,GETDATE()),2),
                                                                                            ('VIP30','percentage',30,DATEADD(DAY,90,GETDATE()),2);

-- =========================================================
-- DISCOUNT USAGE
-- =========================================================
INSERT INTO discount_usage (discount_id,user_id,used_at) VALUES
                                                             (1,3,DATEADD(DAY,-10,GETDATE())),
                                                             (2,5,DATEADD(DAY,-5,GETDATE())),
                                                             (3,2,DATEADD(DAY,-2,GETDATE()));

-- =========================================================
-- REVIEWS
-- =========================================================
INSERT INTO reviews (user_id,product_id,order_item_id,rating,comment) VALUES
                                                                          (1,4,1,5,'Excellent book'),
                                                                          (1,5,1,4,'Good content'),
                                                                          (2,7,3,5,'Great phone'),
                                                                          (3,4,4,4,'Very useful'),
                                                                          (4,9,5,5,'Amazing sound quality');

-- =========================================================
-- REQUESTS
-- =========================================================
INSERT INTO requests (user_id,request_type,status) VALUES
                                                       (1,'refund','open'),
                                                       (2,'support','in_progress'),
                                                       (3,'return','resolved'),
                                                       (4,'support','open'),
                                                       (5,'complaint','rejected');

-- =========================================================
-- FRAUD REPORTS
-- =========================================================
INSERT INTO fraud_reports (stall_id,reporter_id,description) VALUES
                                                                 (5,2,'Reported for fake products'),
                                                                 (5,6,'Reported for late shipping'),
                                                                 (3,10,'Reported for poor quality');