-- name: FullSearch :many
SELECT book_names,
    price,
    author,
    publisher,
    rating,
    genres,
    subgenres,
    ts_rank(searchs_tsv, plainto_tsquery(sqlc.arg(text))) AS ts_rank
FROM searchs
WHERE searchs_tsv @@ plainto_tsquery(sqlc.arg(text))
    AND price BETWEEN sqlc.arg(min_price)
    AND sqlc.arg(max_price)
    AND rating >= sqlc.arg(rating)
ORDER BY ts_rank DESC
LIMIT $1
OFFSET $2;