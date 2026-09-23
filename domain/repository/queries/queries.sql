-- name: CreateTeam :one
INSERT INTO teams (
  name, logo, founded_year, address, city
) VALUES (
  @name, @logo, @founded_year, @address, @city
) RETURNING id;

-- name: GetListTeams :many
SELECT id, name, logo, founded_year
FROM teams
ORDER BY founded_year DESC;

-- name: GetTeamDetail :one
SELECT * FROM teams
WHERE id = @id;

-- name: CheckTeamExisits :one
SELECT id from teams WHERE id = @id;

-- name: UpdateTeams :one
UPDATE teams
SET 
  name = COALESCE(sqlc.narg('name'), name),
  logo = COALESCE(sqlc.narg('logo'), logo),
  founded_year = COALESCE(sqlc.narg('founded_year'), founded_year),
  address = COALESCE(sqlc.narg('address'), address),
  city = COALESCE(sqlc.narg('city'), city),
  updated_at = now()
WHERE id = @id
RETURNING *;

-- name: SoftDeleteTeams :one
UPDATE teams
SET is_deleted = true,
  deleted_at = now()
WHERE id = @id
RETURNING id;

-- name: CreatePlayerTeam :one
INSERT INTO players (
  team_id, name, height_cm, weight_kg, position, jersey_number
) VALUES (
  @team_id, @name, @height_cm, @weight_kg, @position, @jersey_number
) RETURNING id;

-- name: GetListPlayerTeam :many
SELECT id, team_id, name, position, jersey_number
FROM players
WHERE team_id = @team_id
ORDER BY jersey_number ASC;