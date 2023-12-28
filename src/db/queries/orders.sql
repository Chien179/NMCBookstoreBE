-- name: GetOrder :one
SELECT *
FROM orders
WHERE id = $1
LIMIT 1;
-- name: GetOrderToPayment :one
SELECT *
FROM orders
WHERE username = $1
LIMIT 1;
-- name: CreateOrder :one
INSERT INTO orders (username, created_at)
VALUES ($1, NOW())
RETURNING *;
-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;
-- name: ListOdersByUserName :many
SELECT *
FROM orders
WHERE username = $1
ORDER BY id;
-- name: ListOders :one
SELECT t.total_page,
  JSON_AGG(
    json_build_object (
      'id',
      t.id,
      'username',
      t.username,
      'status',
      t.status,
      'sub_amount',
      t.sub_amount,
      'sub_total',
      t.sub_total,
      'sale',
      t.sale,
      'created_at',
      t.created_at
    )
  ) AS orders
FROM (
    SELECT CEILING(
        CAST(COUNT(id) OVER () AS FLOAT) / sqlc.arg('limit')
      ) AS total_page,
      *
    FROM orders
    ORDER BY id
    LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset')
  ) AS t
GROUP BY t.total_page;
-- name: UpdateOrder :one
UPDATE orders
SET status = COALESCE(sqlc.narg(status), status),
  sub_amount = COALESCE(sqlc.narg(sub_amount), sub_amount),
  sub_total = COALESCE(sqlc.narg(sub_total), sub_total),
  sale = COALESCE(sqlc.narg(sale), sale),
  note = COALESCE(sqlc.narg(note), note)
WHERE id = sqlc.arg(id)
RETURNING *;
-- name: ListAllOders :many
SELECT *
FROM orders
ORDER BY id;