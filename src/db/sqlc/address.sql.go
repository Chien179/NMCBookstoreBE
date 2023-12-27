// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: address.sql

package db

import (
	"context"
	"database/sql"
)

const createAddress = `-- name: CreateAddress :one
INSERT INTO address (
    username,
    address,
    district_id,
    city_id
  )
VALUES ($1, $2, $3, $4)
RETURNING id, address, username, city_id, district_id, created_at
`

type CreateAddressParams struct {
	Username   string `json:"username"`
	Address    string `json:"address"`
	DistrictID int64  `json:"district_id"`
	CityID     int64  `json:"city_id"`
}

func (q *Queries) CreateAddress(ctx context.Context, arg CreateAddressParams) (Address, error) {
	row := q.db.QueryRowContext(ctx, createAddress,
		arg.Username,
		arg.Address,
		arg.DistrictID,
		arg.CityID,
	)
	var i Address
	err := row.Scan(
		&i.ID,
		&i.Address,
		&i.Username,
		&i.CityID,
		&i.DistrictID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAddress = `-- name: DeleteAddress :exec
DELETE FROM address
WHERE id = $1
`

func (q *Queries) DeleteAddress(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAddress, id)
	return err
}

const getAddress = `-- name: GetAddress :one
SELECT id, address, username, city_id, district_id, created_at
FROM address
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetAddress(ctx context.Context, id int64) (Address, error) {
	row := q.db.QueryRowContext(ctx, getAddress, id)
	var i Address
	err := row.Scan(
		&i.ID,
		&i.Address,
		&i.Username,
		&i.CityID,
		&i.DistrictID,
		&i.CreatedAt,
	)
	return i, err
}

const listAddresses = `-- name: ListAddresses :many
SELECT address.id AS id,
  address.address AS address,
  districts.name AS district,
  cities.name AS city
FROM address
  INNER JOIN cities ON cities.id = address.city_id
  INNER JOIN districts ON districts.id = address.district_id
WHERE address.username = $1
ORDER BY address.id
`

type ListAddressesRow struct {
	ID       int64  `json:"id"`
	Address  string `json:"address"`
	District string `json:"district"`
	City     string `json:"city"`
}

func (q *Queries) ListAddresses(ctx context.Context, username string) ([]ListAddressesRow, error) {
	rows, err := q.db.QueryContext(ctx, listAddresses, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListAddressesRow{}
	for rows.Next() {
		var i ListAddressesRow
		if err := rows.Scan(
			&i.ID,
			&i.Address,
			&i.District,
			&i.City,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAddress = `-- name: UpdateAddress :one
UPDATE address
SET address = COALESCE($1, address),
  district_id = COALESCE($2, district_id),
  city_id = COALESCE($3, city_id)
WHERE id = $4
RETURNING id, address, username, city_id, district_id, created_at
`

type UpdateAddressParams struct {
	Address    sql.NullString `json:"address"`
	DistrictID sql.NullInt64  `json:"district_id"`
	CityID     sql.NullInt64  `json:"city_id"`
	ID         int64          `json:"id"`
}

func (q *Queries) UpdateAddress(ctx context.Context, arg UpdateAddressParams) (Address, error) {
	row := q.db.QueryRowContext(ctx, updateAddress,
		arg.Address,
		arg.DistrictID,
		arg.CityID,
		arg.ID,
	)
	var i Address
	err := row.Scan(
		&i.ID,
		&i.Address,
		&i.Username,
		&i.CityID,
		&i.DistrictID,
		&i.CreatedAt,
	)
	return i, err
}