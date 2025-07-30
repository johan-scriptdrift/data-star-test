-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = @email;

-- name: GetUsers :many
SELECT
    *
FROM
    users
ORDER BY
    first_name;

-- name: GetLocations :many
SELECT
    *
FROM
    locations
ORDER BY
    created_at;

-- name: GetLocationByID :one
SELECT
    *
FROM
    locations
WHERE id = @id;