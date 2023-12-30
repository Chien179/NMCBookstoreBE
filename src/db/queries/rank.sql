-- name: GetRank :one
SELECT name, score
FROM rank
WHERE  $1 <= score
ORDER BY score DESC
LIMIT 1;