-- name: FullSearch :one
SELECT
    (SELECT (COUNT(*)/sqlc.arg('limit'))
        FROM searchs) 
        as total_page, 
    (SELECT JSON_AGG(t.*) FROM (
        SELECT 
            id,
            name,
            price,
            image,
            description,
            author,
            publisher,
            quantity,
            rating,
            created_at
        FROM searchs
        WHERE searchs_tsv @@ plainto_tsquery(sqlc.arg(text))
            AND searchs.price BETWEEN sqlc.arg(min_price)
            AND sqlc.arg(max_price)
            AND searchs.rating >= sqlc.arg(rating)
        ORDER BY ts_rank DESC
        LIMIT sqlc.arg('limit')
        OFFSET sqlc.arg('offset')
    ) AS t) AS books;