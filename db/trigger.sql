-- From db/migrations/20260314202002_triggers.up.sql
-- Trigger 1: prevent adding items to a locked cart
CREATE TRIGGER trg_prevent_locked_cart ON cart_items
    INSTEAD OF INSERT, UPDATE
                           AS
BEGIN
  IF EXISTS (SELECT 1 FROM inserted i INNER JOIN carts c ON c.cart_id = i.cart_id WHERE c.locked = 1)
BEGIN
    THROW 50001, 'Cannot modify a locked cart', 1;
    RETURN;
END
  -- Proceed with insert/update (SQL Server handles this automatically for INSTEAD OF)
END;

-- Trigger 2: prevent discount code from being used more than once per user
CREATE TRIGGER trg_discount_once ON discount_usage
    INSTEAD OF INSERT
AS
BEGIN
  IF EXISTS (SELECT 1 FROM inserted i WHERE EXISTS (SELECT 1 FROM discount_usage du WHERE du.discount_id = i.discount_id AND du.user_id = i.user_id))
BEGIN
    THROW 50002, 'Discount code already used by user', 1;
    RETURN;
END
END;

-- Trigger 3: auto-suspend stall after 3 fraud reports
CREATE TRIGGER trg_auto_suspend ON fraud_reports
    AFTER INSERT
AS
BEGIN
UPDATE s
SET status = 'suspended'
    FROM stalls s
  INNER JOIN (
    SELECT stall_id
    FROM fraud_reports
    GROUP BY stall_id
    HAVING COUNT(*) >= 3
  ) recent_reports ON s.stall_id = recent_reports.stall_id;
END;

-- Trigger 4: prevent user balance from going negative
CREATE TRIGGER trg_balance_non_negative ON users
    AFTER UPDATE
              AS
BEGIN
  IF EXISTS (SELECT 1 FROM inserted WHERE balance < 0)
BEGIN
    THROW 50003, 'User balance cannot go negative', 1;
END
END;

-- From db/migrations/20260314202033_views.up.sql
CREATE VIEW top_selling_products AS
SELECT TOP 100 PERCENT
    p.product_id,
    p.name,
       p.price,
       p.stall_id,
       CAST(SUM(oi.quantity) AS INT) AS total_sold
FROM order_items oi
         INNER JOIN products p ON p.product_id = oi.product_id
         INNER JOIN orders o ON o.order_id = oi.order_id
WHERE o.order_date >= DATEADD(MONTH, -1, GETDATE())
GROUP BY p.product_id, p.name, p.price, p.stall_id
ORDER BY total_sold DESC;