-- name: FullSearch :one
SELECT t.total_page, JSON_AGG(json_build_object
	('id',id,
	'name',name,
    'price',price,
    'image',image,
    'description',description,
    'author',author,
    'publisher',publisher,
    'quantity',quantity,
    'rating',rating,
    'created_at',created_at)
	) AS books FROM (
        SELECT 
        	CEILING(CAST(COUNT(id) OVER () AS FLOAT)/sqlc.arg('limit')) AS total_page,
            id,
            name,
            price,
            image,
            description,
            author,
            publisher,
            quantity,
            rating,
            created_at,
            ts_rank(searchs_tsv, plainto_tsquery(sqlc.narg(text))) as ts_rank
        FROM searchs
        WHERE
            searchs.price BETWEEN sqlc.arg(min_price) AND sqlc.arg(max_price)
            AND (searchs_tsv @@ plainto_tsquery(sqlc.narg(text)) OR sqlc.narg(text) IS NULL)
            AND (searchs.rating >= sqlc.narg(rating) OR sqlc.narg(rating) IS NULL)
            AND (searchs.genres_id = sqlc.narg(genres_id) OR sqlc.narg(genres_id) IS NULL)
            AND (searchs.subgenres_id = sqlc.narg(subgenres_id) OR sqlc.narg(subgenres_id) IS NULL)
        ORDER BY ts_rank DESC
        LIMIT sqlc.arg('limit')
        OFFSET sqlc.arg('offset')
        ) AS t
    GROUP BY t.total_page;