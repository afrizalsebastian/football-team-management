-- name: CreateTeam :one
INSERT INTO teams (
  name, logo, founded_year, address, city
) VALUES (
  @name, @logo, @founded_year, @address, @city
) RETURNING id;