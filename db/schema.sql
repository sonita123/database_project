

-- =========================================================
-- USERS
-- =========================================================
CREATE TABLE users (
                       user_id    INT IDENTITY(1,1) PRIMARY KEY,
                       first_name NVARCHAR(100) NOT NULL,
                       last_name  NVARCHAR(100) NOT NULL,
                       email      NVARCHAR(255) NOT NULL UNIQUE,
                       password   NVARCHAR(255),
                       balance    DECIMAL(15,2) NOT NULL DEFAULT 0,
                       user_type  NVARCHAR(20) NOT NULL DEFAULT 'regular',
                       created_at DATETIME2(3) NOT NULL DEFAULT GETDATE(),

                       CONSTRAINT chk_users_balance CHECK (balance >= 0),
                       CONSTRAINT chk_users_type CHECK (user_type IN ('regular','vip', 'seller'))
);

-- =========================================================
-- SUPPORTERS
-- =========================================================
CREATE TABLE supporters (
                            supporter_id INT IDENTITY(1,1) PRIMARY KEY,
                            first_name   NVARCHAR(100) NOT NULL,
                            last_name    NVARCHAR(100) NOT NULL,
                            image_url    NVARCHAR(MAX),
                            email        NVARCHAR(255),
                            password     NVARCHAR(255)
);

-- =========================================================
-- ADMINS
-- =========================================================
CREATE TABLE admins (
                        admin_id   INT IDENTITY(1,1) PRIMARY KEY,
                        username   NVARCHAR(100) NOT NULL UNIQUE,
                        password   NVARCHAR(255) NOT NULL,
                        created_at DATETIME2(3) NOT NULL DEFAULT GETDATE()
);

-- =========================================================
-- SESSIONS (FIXED 🔥)
-- =========================================================
CREATE TABLE sessions (
                          token CHAR(64) PRIMARY KEY,  -- ✅ FIXED (was NVARCHAR(500))
                          user_id INT NOT NULL,
                          role NVARCHAR(50) NOT NULL,
                          expires_at DATETIME2(3) NOT NULL,

                          CONSTRAINT fk_sessions_user FOREIGN KEY (user_id)
                              REFERENCES users(user_id) ON DELETE CASCADE
);

-- =========================================================
-- VIP USERS
-- =========================================================
CREATE TABLE vip_users (
                           vip_id     INT IDENTITY(1,1) PRIMARY KEY,
                           user_id    INT NOT NULL,
                           start_date DATETIME2(3) NOT NULL DEFAULT GETDATE(),
                           end_date   DATETIME2(3) NOT NULL,

                           CONSTRAINT fk_vip_user FOREIGN KEY (user_id)
                               REFERENCES users(user_id) ON DELETE CASCADE
);

-- =========================================================
-- ADDRESSES
-- =========================================================
CREATE TABLE addresses (
                           address_id  INT IDENTITY(1,1) PRIMARY KEY,
                           user_id     INT NOT NULL,
                           city        NVARCHAR(100) NOT NULL,
                           street      NVARCHAR(255) NOT NULL,
                           postal_code NVARCHAR(20),
                           is_default  BIT NOT NULL DEFAULT 0,

                           CONSTRAINT fk_address_user FOREIGN KEY (user_id)
                               REFERENCES users(user_id) ON DELETE CASCADE
);

-- =========================================================
-- SELLERS
-- =========================================================
CREATE TABLE sellers (
                         seller_id     INT IDENTITY(1,1) PRIMARY KEY,
                         user_id       INT NOT NULL,
                         registered_at DATETIME2(3) NOT NULL DEFAULT GETDATE(),

                         CONSTRAINT fk_seller_user FOREIGN KEY (user_id)
                             REFERENCES users(user_id) ON DELETE CASCADE
);

-- =========================================================
-- STALLS
-- =========================================================
CREATE TABLE stalls (
                        stall_id    INT IDENTITY(1,1) PRIMARY KEY,
                        seller_id   INT NOT NULL,
                        name        NVARCHAR(255) NOT NULL,
                        status      NVARCHAR(20) NOT NULL DEFAULT 'pending',
                        approved_by INT,
                        approved_at DATETIME2(3),
                        created_at  DATETIME2(3) NOT NULL DEFAULT GETDATE(),

                        CONSTRAINT chk_stall_status CHECK (
                            status IN ('pending','active','suspended','inactive')
                            ),

                        CONSTRAINT fk_stall_seller FOREIGN KEY (seller_id)
                            REFERENCES sellers(seller_id) ON DELETE CASCADE,

                        CONSTRAINT fk_stall_supporter FOREIGN KEY (approved_by)
                            REFERENCES supporters(supporter_id) ON DELETE SET NULL
);

-- =========================================================
-- PRODUCTS
-- =========================================================
CREATE TABLE products (
                          product_id INT IDENTITY(1,1) PRIMARY KEY,
                          stall_id   INT NOT NULL,
                          name       NVARCHAR(255) NOT NULL,
                          price      DECIMAL(15,2) NOT NULL,
                          stock      INT NOT NULL DEFAULT 0,
                          status     NVARCHAR(20) NOT NULL DEFAULT 'active',
                          created_at DATETIME2(3) NOT NULL DEFAULT GETDATE(),

                          CONSTRAINT chk_product_price CHECK (price >= 0),
                          CONSTRAINT chk_product_stock CHECK (stock >= 0),
                          CONSTRAINT chk_product_status CHECK (
                              status IN ('active','inactive','out_of_stock','locked')
                              ),

                          CONSTRAINT fk_product_stall FOREIGN KEY (stall_id)
                              REFERENCES stalls(stall_id) ON DELETE CASCADE
);

-- =========================================================
-- PRODUCT VIEWS (FIXED 🔥)
-- =========================================================
CREATE TABLE product_views (
                               view_id    INT IDENTITY(1,1) PRIMARY KEY,
                               user_id    INT NOT NULL,
                               product_id INT NOT NULL,
                               viewed_at  DATETIME2(3) NOT NULL DEFAULT GETDATE(),

                               CONSTRAINT fk_view_user FOREIGN KEY (user_id)
                                   REFERENCES users(user_id) ON DELETE CASCADE,

    -- ✅ FIXED (removed cascade conflict)
                               CONSTRAINT fk_view_product FOREIGN KEY (product_id)
                                   REFERENCES products(product_id) ON DELETE NO ACTION
);

-- =========================================================
-- REQUESTS
-- =========================================================
CREATE TABLE requests (
                          request_id   INT IDENTITY(1,1) PRIMARY KEY,
                          user_id      INT NOT NULL,
                          request_type NVARCHAR(50) NOT NULL,
                          status       NVARCHAR(20) NOT NULL DEFAULT 'open',
                          description  NVARCHAR(MAX),
                          handled_by   INT,
                          resolved_at  DATETIME2(3),
                          created_at   DATETIME2(3) NOT NULL DEFAULT GETDATE(),

                          CONSTRAINT chk_request_status CHECK (
                              status IN ('open','in_progress','resolved','rejected')
                              ),

                          CONSTRAINT fk_request_user FOREIGN KEY (user_id)
                              REFERENCES users(user_id) ON DELETE CASCADE,

                          CONSTRAINT fk_request_supporter FOREIGN KEY (handled_by)
                              REFERENCES supporters(supporter_id) ON DELETE SET NULL
);

-- =========================================================
-- CARTS
-- =========================================================
CREATE TABLE carts (
                       cart_id    INT IDENTITY(1,1) PRIMARY KEY,
                       user_id    INT NOT NULL,
                       locked     BIT NOT NULL DEFAULT 0,
                       created_at DATETIME2(3) NOT NULL DEFAULT GETDATE(),

                       CONSTRAINT fk_cart_user FOREIGN KEY (user_id)
                           REFERENCES users(user_id) ON DELETE CASCADE
);

-- =========================================================
-- CART ITEMS
-- =========================================================
CREATE TABLE cart_items (
                            cart_item_id INT IDENTITY(1,1) PRIMARY KEY,
                            cart_id      INT NOT NULL,
                            product_id   INT NOT NULL,
                            quantity     INT NOT NULL DEFAULT 1,

                            CONSTRAINT chk_cart_qty CHECK (quantity > 0),

                            CONSTRAINT fk_cart_items_cart FOREIGN KEY (cart_id)
                                REFERENCES carts(cart_id) ON DELETE CASCADE,
                            CONSTRAINT fk_cart_items_product FOREIGN KEY (product_id)
                                REFERENCES products(product_id) ON DELETE NO ACTION
);

-- =========================================================
-- ORDERS (FIXED 🔥)
-- =========================================================
CREATE TABLE orders (
                        order_id    INT IDENTITY(1,1) PRIMARY KEY,
                        user_id     INT NOT NULL,
                        cart_id     INT NULL,
                        address_id  INT NULL,
                        total_price DECIMAL(15,2) NOT NULL,
                        status      NVARCHAR(20) NOT NULL DEFAULT 'pending',
                        order_date  DATETIME2(3) NOT NULL DEFAULT GETDATE(),

                        CONSTRAINT chk_order_total CHECK (total_price >= 0),
                        CONSTRAINT chk_order_status CHECK (
                            status IN ('pending','confirmed','shipped','delivered','cancelled')
                            ),
                        CONSTRAINT fk_address_users FOREIGN KEY (user_id)
                            REFERENCES users(user_id) ON DELETE NO ACTION,

    -- ✅ FIXED
                        CONSTRAINT fk_order_cart FOREIGN KEY (cart_id)
                            REFERENCES carts(cart_id) ON DELETE NO ACTION,

                        CONSTRAINT fk_order_address FOREIGN KEY (address_id)
                            REFERENCES addresses(address_id) ON DELETE SET NULL
);

-- =========================================================
-- ORDER ITEMS
-- =========================================================
CREATE TABLE order_items (
                             order_item_id INT IDENTITY(1,1) PRIMARY KEY,
                             order_id      INT NOT NULL,
                             product_id    INT NOT NULL,
                             quantity      INT NOT NULL DEFAULT 1,
                             price         DECIMAL(15,2) NOT NULL,

                             CONSTRAINT chk_order_item_qty CHECK (quantity > 0),
                             CONSTRAINT chk_order_item_price CHECK (price >= 0),

                             CONSTRAINT fk_order_items_order FOREIGN KEY (order_id)
                                 REFERENCES orders(order_id) ON DELETE CASCADE,

                             CONSTRAINT fk_order_items_product FOREIGN KEY (product_id)
                                 REFERENCES products(product_id) ON DELETE CASCADE
);

-- =========================================================
-- DISCOUNT CODES
-- =========================================================
CREATE TABLE discount_codes (
                                discount_id     INT IDENTITY(1,1) PRIMARY KEY,
                                code            NVARCHAR(50) NOT NULL UNIQUE,
                                discount_type   NVARCHAR(20) NOT NULL DEFAULT 'percentage',
                                percentage      DECIMAL(5,2),
                                fixed_amount    DECIMAL(15,2),
                                expiration_date DATETIME2(3) NOT NULL,
                                max_uses        INT,
                                is_active       BIT NOT NULL DEFAULT 1,
                                supporter_id    INT,

                                CONSTRAINT chk_discount_type CHECK (discount_type IN ('percentage','fixed')),
                                CONSTRAINT chk_percentage CHECK (percentage IS NULL OR (percentage > 0 AND percentage <= 100)),
                                CONSTRAINT chk_fixed_amount CHECK (fixed_amount IS NULL OR fixed_amount > 0),

                                CONSTRAINT fk_discount_supporter FOREIGN KEY (supporter_id)
                                    REFERENCES supporters(supporter_id) ON DELETE SET NULL
);

-- =========================================================
-- DISCOUNT USAGE
-- =========================================================
CREATE TABLE discount_usage (
                                usage_id    INT IDENTITY(1,1) PRIMARY KEY,
                                discount_id INT NOT NULL,
                                user_id     INT NOT NULL,
                                used_at     DATETIME2(3) NOT NULL DEFAULT GETDATE(),

                                CONSTRAINT uq_discount_user UNIQUE (discount_id, user_id),

                                CONSTRAINT fk_usage_discount FOREIGN KEY (discount_id)
                                    REFERENCES discount_codes(discount_id) ON DELETE CASCADE,

                                CONSTRAINT fk_usage_user FOREIGN KEY (user_id)
                                    REFERENCES users(user_id) ON DELETE CASCADE
);

-- =========================================================
-- REVIEWS
-- =========================================================
CREATE TABLE reviews (
                         review_id     INT IDENTITY(1,1) PRIMARY KEY,
                         user_id       INT NOT NULL,
                         product_id    INT NOT NULL,
                         order_item_id INT NULL,
                         rating        INT NOT NULL,
                         comment       NVARCHAR(MAX),
                         created_at    DATETIME2(3) NOT NULL DEFAULT GETDATE(),

                         CONSTRAINT chk_rating CHECK (rating BETWEEN 1 AND 5),
                         CONSTRAINT uq_review UNIQUE (user_id, product_id),

    -- ✅ ownership
                         CONSTRAINT fk_review_user FOREIGN KEY (user_id)
                             REFERENCES users(user_id) ON DELETE CASCADE,

    -- ✅ already fixed earlier
                         CONSTRAINT fk_review_product FOREIGN KEY (product_id)
                             REFERENCES products(product_id) ON DELETE NO ACTION,

    -- ❗ FINAL FIX HERE
                         CONSTRAINT fk_review_order_item FOREIGN KEY (order_item_id)
                             REFERENCES order_items(order_item_id) ON DELETE NO ACTION
);

-- =========================================================
-- FRAUD REPORTS
-- =========================================================

CREATE TABLE fraud_reports (
                               report_id   INT IDENTITY(1,1) PRIMARY KEY,
                               stall_id    INT NOT NULL,
                               reporter_id INT NOT NULL,
                               description NVARCHAR(MAX),
                               reported_at DATETIME2(3) NOT NULL DEFAULT GETDATE(),

                               CONSTRAINT fk_report_stall FOREIGN KEY (stall_id)
                                   REFERENCES stalls(stall_id) ON DELETE CASCADE,

    -- ❗ FIXED: prevent cascade conflict
                               CONSTRAINT fk_report_user FOREIGN KEY (reporter_id)
                                   REFERENCES users(user_id) ON DELETE NO ACTION
);