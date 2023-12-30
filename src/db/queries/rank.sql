-- name: GetRank :one
SELECT name, score
FROM rank
WHERE $1 <= score
ORDER BY score ASC
LIMIT 1;