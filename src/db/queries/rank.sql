-- name: GetRank :one
SELECT name, score
FROM rank
WHERE score <= $1
ORDER BY score DESC
LIMIT 1;