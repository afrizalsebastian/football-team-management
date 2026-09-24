-- name: CreateTeam :one
INSERT INTO teams (
  name, logo, founded_year, address, city
) VALUES (
  @name, @logo, @founded_year, @address, @city
) RETURNING id;

-- name: GetListTeams :many
SELECT id, name, logo, founded_year
FROM teams
WHERE is_deleted = false
ORDER BY founded_year DESC;

-- name: GetTeamDetail :one
SELECT * FROM teams
WHERE id = @id AND is_deleted = false;

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
WHERE id = @id AND is_deleted = false
RETURNING *;

-- name: SoftDeleteTeams :one
UPDATE teams
SET is_deleted = true,
  deleted_at = now()
WHERE id = @id AND is_deleted = false
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
WHERE team_id = @team_id AND is_deleted = false
ORDER BY jersey_number ASC;

-- name: GetListPlayer :many
SELECT 
  p.id as id,
  p.name,
  p.position,
  p.jersey_number,
  t.id as team_id,
  t.name as team_name
FROM players p
LEFT JOIN teams t ON p.team_id = t.id AND t.is_deleted = false
WHERE p.is_deleted = false
ORDER BY p.created_at ASC;

-- name: UpdatePlayers :one
UPDATE players
SET 
  team_id = COALESCE(sqlc.narg('team_id'), team_id),
  name = COALESCE(sqlc.narg('name'), name),
  height_cm = COALESCE(sqlc.narg('height_cm'), height_cm),
  weight_kg = COALESCE(sqlc.narg('weight_kg'), weight_kg),
  position = COALESCE(sqlc.narg('position'), position),
  jersey_number = COALESCE(sqlc.narg('jersey_number'), jersey_number),
  updated_at = now()
WHERE id = @id AND is_deleted = false
RETURNING *;

-- name: CreateMatches :one
INSERT INTO matches(
  match_date, match_time, home_team_id, away_team_id
) VALUES (
  @match_date, @match_time, @home_team_id, @away_team_id
) RETURNING id;

-- name: GetListMatch :many
SELECT
  m.id,
  m.match_date,
  m.match_time,
  m.home_team_id,
  m.away_team_id,
  home.name as home_team_name,
  home.logo as home_logo,
  away.name as away_team_name,
  away.logo as away_logo
FROM matches m
LEFT JOIN teams home ON m.home_team_id = home.id AND home.is_deleted = false
LEFT JOIN teams away ON m.away_team_id = away.id AND away.is_deleted = false
WHERE m.is_deleted = false
  AND (
    (sqlc.narg('home_team_id')::uuid IS NULL OR home_team_id = sqlc.narg('home_team_id')::uuid) OR 
    (sqlc.narg('away_team_id')::uuid IS NULL OR away_team_id = sqlc.narg('away_team_id')::uuid)
  )
ORDER BY m.match_date ASC, m.match_time ASC;

-- name: RescheduleMatch :one
UPDATE matches
SET 
  match_date = @match_date,
  match_time = @match_time,
  updated_at = now()
WHERE id = @id and is_deleted = false
RETURNING id;

-- name: DeleteMatch :one
UPDATE matches
SET 
  is_deleted = true,
  deleted_at = now()
WHERE id = @id and is_deleted = false
RETURNING id;