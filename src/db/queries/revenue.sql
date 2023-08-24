
-- name: RevenueDays :many
SELECT
COALESCE(SUM(subtotal), 0) AS sum_revenue,
COALESCE(AVG(subtotal), 0) AS avg_revenue,  
to_char(date(created_at),'DD-MM-YYYY') as time_revenue,
to_char(date(created_at),'YYYY') as time_year
FROM payments
WHERE status = 'success'
GROUP BY time_revenue, time_year
ORDER BY time_year ASC;

-- name: RevenueMonths :many
SELECT
COALESCE(SUM(subtotal), 0) AS sum_revenue,
COALESCE(AVG(subtotal), 0) AS avg_revenue,
to_char(date(created_at),'MM-YYYY') as time_revenue,
to_char(date(created_at),'YYYY') as time_year
FROM payments
WHERE status = 'success'
GROUP BY time_revenue, time_year
ORDER BY time_year ASC;

-- name: RevenueQuarters :many
SELECT
COALESCE(SUM(subtotal), 0) AS sum_revenue,
COALESCE(AVG(subtotal), 0) AS avg_revenue,
to_char(date(created_at),'Q-YYYY') as time_revenue,
to_char(date(created_at),'YYYY') as time_year
FROM payments
WHERE status = 'success'
GROUP BY time_revenue, time_year
ORDER BY time_year ASC;

-- name: RevenueYears :many
SELECT
COALESCE(SUM(subtotal), 0) AS sum_revenue,
COALESCE(AVG(subtotal), 0) AS avg_revenue,
to_char(date(created_at),'YYYY') as time_revenue
FROM payments
WHERE status = 'success'
GROUP BY time_revenue
ORDER BY time_revenue;