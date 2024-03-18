SELECT category, COUNT(*) AS total_products
FROM products
GROUP BY category;

SELECT
SUM(CASE WHEN stock BETWEEN 1 AND 100 THEN 1 ELSE 0 END) AS stock_1_to_100,
SUM(CASE WHEN stock BETWEEN 101 AND 200 THEN 1 ELSE 0 END) AS stock_101_to_200,
SUM(CASE WHEN stock BETWEEN 201 AND 300 THEN 1 ELSE 0 END) AS stock_201_to_300
FROM inventory;


SELECT
CASE
WHEN stock BETWEEN 1 AND 100 THEN '1-100'
WHEN stock BETWEEN 101 AND 200 THEN '101-200'
WHEN stock BETWEEN 201 AND 300 THEN '201-300'
ELSE '其他'
END AS stock_range,
GROUP_CONCAT(product_name) AS product_names
FROM products
Left JOIN inventory ON products.product_id = inventory.product_id
GROUP BY
CASE
WHEN stock BETWEEN 1 AND 100 THEN '1-100'
WHEN stock BETWEEN 101 AND 200 THEN '101-200'
WHEN stock BETWEEN 201 AND 300 THEN '201-300'
ELSE '其他'
END;

用户表（users）：
user_id (主键)
username
email
其他用户信息字段...
订单表（orders）：
order_id (主键)
user_id (外键，关联用户表的 user_id)
total_amount (订单总金额)
订单其他信息字段...
商品表（products）：
product_id (主键)
product_name
price (商品单价)
商品其他信息字段...
订单和商品的中间表（order_products）：
order_id (外键，关联订单表的 order_id)
product_id (外键，关联商品表的 product_id)
quantity (购买数量)
其他中间表字段...


SELECT o.order_id, o.total_amount, op.product_id, p.product_name, op.quantity
FROM orders o
JOIN order_products op ON o.order_id = op.order_id
JOIN products p ON op.product_id = p.product_id
WHERE o.user_id = [用户ID];


SELECT COUNT(*) AS total_orders
FROM orders;


SELECT u.user_id, SUM(o.total_amount) AS total_payment
FROM users u
JOIN orders o ON u.user_id = o.user_id
GROUP BY u.user_id;