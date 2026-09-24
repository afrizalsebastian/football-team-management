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
SELECT id FROM teams WHERE id = @id;

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
LEFT JOIN teams t ON p.team_id = t.id
WHERE p.is_deleted = false
  AND t.is_deleted = false
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

-- name: GetPlayerDetail :one
SELECT 
  p.id as id,
  p.name,
  p.position,
  p.jersey_number,
  p.height_cm,
  p.weight_kg,
  t.id as team_id,
  t.name as team_name,
  count(g.id) as goals_count,
  p.created_at,
  p.updated_at,
  p.is_deleted,
  p.deleted_at
FROM players p
LEFT JOIN teams t ON p.team_id = t.id
LEFT JOIN goals g ON g.player_id = p.id AND g.is_deleted = false
WHERE p.id = @id
  AND p.is_deleted = false
  AND t.is_deleted = false
GROUP BY p.id, t.id
ORDER BY p.created_at ASC;

-- name: DeletePlayer :one
UPDATE players
SET 
  is_deleted = true,
  deleted_at = now()
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
  away.logo as away_logo,
  mr.status,
  mr.home_score,
  mr.away_score,
  p.id as player_id,
  p.name as player_name,
  p.jersey_number as player_number,
  p.team_id as player_team_id,
  g.id as goal_id,
  g.goal_minute
FROM matches m
LEFT JOIN teams home ON m.home_team_id = home.id
LEFT JOIN teams away ON m.away_team_id = away.id
LEFT JOIN match_results mr ON mr.id = m.id
LEFT JOIN goals g ON g.match_id = m.id and g.is_deleted = false
LEFT JOIN players p ON p.id = g.player_id
WHERE m.is_deleted = false
  AND (
    (sqlc.narg('home_team_id')::uuid IS NULL OR home_team_id = sqlc.narg('home_team_id')::uuid) OR 
    (sqlc.narg('away_team_id')::uuid IS NULL OR away_team_id = sqlc.narg('away_team_id')::uuid)
  )
  AND home.is_deleted = false
  AND away.is_deleted = false
  AND (
    sqlc.narg('filter_type')::text IS NULL OR
    (sqlc.narg('filter_type')::text = 'scheduled' AND mr.id IS NULL) OR
    (sqlc.narg('filter_type')::text = 'result' AND mr.id IS NOT NULL)
  )
ORDER BY m.match_date ASC, m.match_time ASC, g.goal_minute ASC;

-- name: RescheduleMatch :one
UPDATE matches
SET 
  match_date = @match_date,
  match_time = @match_time,
  updated_at = now()
WHERE id = @id AND is_deleted = false
RETURNING id;

-- name: DeleteMatch :one
UPDATE matches
SET 
  is_deleted = true,
  deleted_at = now()
WHERE id = @id AND is_deleted = false
RETURNING id;

-- name: CreateGoals :one
INSERT INTO goals(
  match_id, player_id, goal_minute
) 
SELECT
  m.id,
  p.id,
  @goal_minute
FROM matches m
JOIN players p ON p.id = @player_id 
WHERE m.id = @match_id
  AND m.is_deleted = false
  AND p.is_deleted = false 
  AND (m.home_team_id = p.team_id OR m.away_team_id = p.team_id)
RETURNING id;

-- name: GetListMatchGoals :many
SELECT
  g.id,
  g.match_id,
  home.id as home_team_id,
  home.name as home_team_name,
  away.id as away_team_id,
  away.name as away_team_name,
  g.player_id,
  p.name as player_name,
  p.jersey_number as jersey_number,
  t.id as team_id,
  t.name as team_name,
  g.goal_minute
FROM goals g
JOIN matches m ON m.id = g.match_id
JOIN players p ON p.id = g.player_id
JOIN teams t ON t.id = p.team_id
LEFT JOIN teams home ON m.home_team_id = home.id
LEFT JOIN teams away ON m.away_team_id = away.id
WHERE g.match_id = @match_id AND g.is_deleted = false
ORDER BY g.created_at ASC;

-- name: DeleteGoals :one
UPDATE goals
SET
  is_deleted = true,
  deleted_at = now()
WHERE id = @id AND is_deleted = false
RETURNING id; 

-- name: MatchFullTime :one
WITH match_info AS (
  SELECT ma.id as match_id, home_team_id, away_team_id
  FROM matches ma
  WHERE ma.id = @match_id
),
goals_count AS (
  SELECT
    m.match_id,
    COALESCE(SUM(CASE WHEN p.team_id = m.home_team_id THEN 1 ELSE 0 END), 0) AS home_score,
    COALESCE(SUM(CASE WHEN p.team_id = m.away_team_id THEN 1 ELSE 0 END), 0) AS away_score
  FROM match_info m
  LEFT JOIN goals g ON g.match_id = m.match_id and g.is_deleted = false
  LEFT JOIN players p ON p.id = g.player_id
  GROUP BY m.match_id
)
INSERT INTO match_results(
  id, status, home_score, away_score
)
SELECT
  gc.match_id as id,
  CASE 
      WHEN gc.home_score > gc.away_score THEN 1 
      WHEN gc.home_score < gc.away_score THEN -1
      ELSE 0                            
  END AS status,
  gc.home_score,
  gc.away_score
FROM goals_count gc
ON CONFLICT (id) 
DO UPDATE SET 
    status = EXCLUDED.status,
    home_score = EXCLUDED.home_score,
    away_score = EXCLUDED.away_score
RETURNING id;