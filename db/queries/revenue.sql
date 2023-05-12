-- name: RevenueHours :many
SELECT
COALESCE(SUM(subtotal), 0) AS revenue_days,
COALESCE(AVG(subtotal), 0) AS avg_revenue_days,  
created_at
FROM payments
WHERE status = 'success'
GROUP BY created_at
ORDER BY created_at;

-- name: RevenueDays :many
SELECT
COALESCE(SUM(subtotal), 0) AS revenue_days,
COALESCE(AVG(subtotal), 0) AS avg_revenue_days,  
to_char(date(created_at),'YYYY-MM-DD') as dates
FROM payments
WHERE status = 'success'
GROUP BY dates
ORDER BY dates;

-- name: RevenueMonths :many
SELECT
COALESCE(SUM(subtotal), 0) AS revenue_months,
COALESCE(AVG(subtotal), 0) AS avg_revenue_months,
to_char(date(created_at),'YYYY-MM') as year_months
FROM payments
WHERE status = 'success'
GROUP BY year_months
ORDER BY year_months;

-- name: RevenueQuarters :many
SELECT
COALESCE(SUM(subtotal), 0) AS revenue_quarters,
COALESCE(AVG(subtotal), 0) AS avg_revenue_quarters,
to_char(date(created_at),'YYYY-Q') as year_quarters
FROM payments
WHERE status = 'success'
GROUP BY year_quarters
ORDER BY year_quarters;

-- name: RevenueYears :many
SELECT
COALESCE(SUM(subtotal), 0) AS revenue_years,
COALESCE(AVG(subtotal), 0) AS avg_revenue_years,
to_char(date(created_at),'YYYY') as years
FROM payments
WHERE status = 'success'
GROUP BY years
ORDER BY years;